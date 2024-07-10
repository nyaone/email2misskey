package handlers

import (
	"bytes"
	"email2misskey/config"
	"email2misskey/consts"
	"email2misskey/global"
	"email2misskey/misskey"
	"fmt"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"strings"
	"time"
)

func IncomingEMAil() backends.ProcessorConstructor {
	return func() backends.Decorator {
		return func(p backends.Processor) backends.Processor {
			return backends.ProcessWith(
				func(e *mail.Envelope, task backends.SelectTask) (backends.Result, error) {
					if task == backends.TaskSaveMail {
						global.Logger.Debugf("A new incoming eMail")
						// Check size
						if config.Config.EMail.SizeLimit > 0 && e.Data.Len() > config.Config.EMail.SizeLimit {
							// Too large
							err := fmt.Errorf("email too large")
							return backends.NewResult(fmt.Sprintf("552 Error: %s", err)), err
						}

						// Verify DKIM
						if config.Config.EMail.VerifyDKIM {
							verifications, err := dkim.Verify(bytes.NewReader(e.Data.Bytes()))
							if err != nil {
								// Failed to verify DKIM
								return backends.NewResult(fmt.Sprintf("451 Error: %s", err)), err
							}
							if len(verifications) == 0 {
								// No signature found
								err = fmt.Errorf("no DKIM signature found")
								return backends.NewResult(fmt.Sprintf("503 Error: %s", err)), err
							}
							for _, v := range verifications {
								if v.Err != nil {
									// Invalid signature
									err = fmt.Errorf("email signature invalid")
									return backends.NewResult(fmt.Sprintf("503 Error: %s", err)), err
								}
							}
						}

						// Send to target user
						var pendingUserIDs []string
						for _, rcpt := range e.RcptTo {
							username := strings.ToLower(rcpt.User)
							userID, err := misskey.GetTargetUserID(username)
							if err != nil {
								// Network failed
								return backends.NewResult(fmt.Sprintf("554 Error: %s", err)), err
							}
							if userID != nil {
								pendingUserIDs = append(pendingUserIDs, *userID)
							}
						}

						if len(pendingUserIDs) > 0 {

							// Upload eMail to Misskey
							filename := fmt.Sprintf(consts.EMailFilenameTemplate, time.Now().Format("2006-0102-150405"))
							fileID, fileURL, err := misskey.UploadFile(filename, &e.Data)
							if err != nil {
								// Network failed
								global.Logger.Errorf("Failed to upload file to Misskey")
								return backends.NewResult(fmt.Sprintf("554 Error: %s", err)), err
							}

							// Create note summary & message
							summary := fmt.Sprintf(consts.SummaryTemplate, e.Subject)
							detail := fmt.Sprintf(
								consts.MessageTemplate,
								e.MailFrom.String(),
								fmt.Sprintf("%s/mail?url=%s", config.Config.EMail.ReaderUrl, fileURL),
							)

							// Send note
							err = misskey.CreatePrivateNote(pendingUserIDs, summary, detail, fileID)
							if err != nil {
								// Network failed
								global.Logger.Errorf("Failed to send message to users %v for email %s, might need to resend later", pendingUserIDs, fileID)
							}

							// All done, pass to next processor
							return p.Process(e, task)
						} else {
							// No match, wrong message
							err := fmt.Errorf("no user found")
							return backends.NewResult(fmt.Sprintf("551 Error: %s", err)), err
						}

					} else {

						// Bypass others
						return p.Process(e, task)
					}
				},
			)
		}
	}
}
