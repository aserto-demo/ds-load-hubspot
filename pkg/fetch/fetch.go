package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aserto-demo/ds-load-hubspot/pkg/hubspotclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
)

type Fetcher struct {
	hubspotClient *hubspotclient.HubspotClient
	Contacts      bool
	Companies     bool
}

func New(ctx context.Context, privateAccessToken, clientID, clientSecret, refreshToken string) (*Fetcher, error) {
	var client *hubspotclient.HubspotClient
	var err error
	if privateAccessToken != "" {
		client, err = hubspotclient.NewHubspotClient(ctx, privateAccessToken)
	} else {
		client, err = hubspotclient.NewHubspotOAuth2Client(ctx, clientID, clientSecret, refreshToken)
	}
	if err != nil {
		return nil, err
	}

	return &Fetcher{
		hubspotClient: client,
	}, nil
}

func (f *Fetcher) WithOptions(contacts, companies bool) *Fetcher {
	f.Contacts = contacts
	f.Companies = companies
	return f
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	writer, err := js.NewJSONArrayWriter(outputWriter)
	if err != nil {
		return err
	}
	defer writer.Close()

	users, err := f.hubspotClient.ListUsers()
	if err != nil {
		_, _ = errorWriter.Write([]byte(err.Error()))
		common.SetExitCode(1)
		return err
	}

	for _, user := range users {
		userBytes, err := json.Marshal(user)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			continue
		}
		var obj map[string]interface{}
		err = json.Unmarshal(userBytes, &obj)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			continue
		}
		obj["type"] = "user"

		err = writer.Write(obj)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}
	}

	if f.Companies {
		companies, err := f.hubspotClient.ListCompanies()
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			return err
		}

		for _, company := range companies {
			companyBytes, err := json.Marshal(company)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				continue
			}
			var obj map[string]interface{}
			err = json.Unmarshal(companyBytes, &obj)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				continue
			}
			obj["type"] = "company"

			err = writer.Write(obj)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
			}
		}
	}

	if f.Contacts {
		contacts, err := f.hubspotClient.ListContacts()
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
			common.SetExitCode(1)
			return err
		}

		for _, contact := range contacts {
			contactBytes, err := json.Marshal(contact)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				continue
			}
			var obj map[string]interface{}
			err = json.Unmarshal(contactBytes, &obj)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
				common.SetExitCode(1)
				continue
			}
			obj["key"] = contact.Properties.Email
			obj["displayName"] = createDisplayName(contact.Properties.FirstName, contact.Properties.LastName, contact.Properties.Email)
			obj["type"] = "contact"
			if contact.Properties.Owner != "" {
				user := f.hubspotClient.LookupUser(contact.Properties.Owner)
				if user != "" {
					obj["owner"] = user
				}
			}

			if contact.Properties.Company != "" {
				company := f.hubspotClient.LookupCompany(contact.Properties.Company)
				if company.ID != "" {
					obj["companyId"] = company.ID
				}
			}

			err = writer.Write(obj)
			if err != nil {
				_, _ = errorWriter.Write([]byte(err.Error()))
			}
		}
	}

	return nil
}

func createDisplayName(firstName, lastName, email string) string {
	returnVal := firstName
	if returnVal != "" {
		if lastName != "" {
			returnVal = fmt.Sprintf("%s %s", firstName, lastName)
		}
	} else {
		if lastName != "" {
			returnVal = lastName
		} else {
			returnVal = email
		}
	}
	return returnVal
}
