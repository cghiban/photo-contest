package utils

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type state_keys []string

func (k state_keys) Len() int {
	return len(k)
}
func (s state_keys) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s state_keys) Less(i, j int) bool {
	if s[i] == "NY" || s[j] == "OO" {
		return true
	}
	if s[i] == "OO" || s[j] == "NY" {
		return false
	}
	return s[i] < s[j]
}

func USStates() map[string]string {
	return map[string]string{
		"NY": "New York",
		"AL": "Alabama",
		"AK": "Alaska",
		"AZ": "Arizona",
		"AR": "Arkansas",
		"CA": "California",
		"CO": "Colorado",
		"CT": "Connecticut",
		"DE": "Delaware",
		"DC": "District Of Columbia",
		"FL": "Florida",
		"GA": "Georgia",
		"HI": "Hawaii",
		"ID": "Idaho",
		"IL": "Illinois",
		"IN": "Indiana",
		"IA": "Iowa",
		"KS": "Kansas",
		"KY": "Kentucky",
		"LA": "Louisiana",
		"ME": "Maine",
		"MD": "Maryland",
		"MA": "Massachusetts",
		"MI": "Michigan",
		"MN": "Minnesota",
		"MS": "Mississippi",
		"MO": "Missouri",
		"MT": "Montana",
		"NE": "Nebraska",
		"NV": "Nevada",
		"NH": "New Hampshire",
		"NJ": "New Jersey",
		"NM": "New Mexico",
		"NC": "North Carolina",
		"ND": "North Dakota",
		"OH": "Ohio",
		"OK": "Oklahoma",
		"OR": "Oregon",
		"PA": "Pennsylvania",
		"RI": "Rhode Island",
		"SC": "South Carolina",
		"SD": "South Dakota",
		"TN": "Tennessee",
		"TX": "Texas",
		"UT": "Utah",
		"VT": "Vermont",
		"VA": "Virginia",
		"WA": "Washington",
		"WV": "West Virginia",
		"WI": "Wisconsin",
		"WY": "Wyoming",
		"OO": "Out of US",
	}
}

func StateKeys(states map[string]string) state_keys {
	i := 0
	keys := make(state_keys, len(states))
	for k := range states {
		keys[i] = k
		i++
	}
	return keys
}

func Ethnicities() map[string]string {
	return map[string]string{
		"as": "Asian/Pacific Islander",
		"aa": "Black or African American",
		"hs": "Hispanic or Latino",
		"na": "Native American or American Indian",
		"oo": "Other",
		"wh": "White",
		"pn": "Prefer not to respond",
	}
}
func EthnicitiesKeys() []string {
	return []string{"as", "aa", "hs", "na", "oo", "wh", "pn"}
}

func Genders() map[string]string {
	return map[string]string{
		"M": "Male",
		"F": "Female",
		"N": "Non-Binary",
		"-": "Prefer not to answer",
	}
}
func GenderKeys() []string {
	return []string{"M", "F", "N", "-"}
}

func InStringSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ReadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer file.Close()

	b, _ := ioutil.ReadAll(file)
	return string(b), nil
}

func SendEmail(mg *mailgun.MailgunImpl, sender, subject, recipient, bodyText, bodyHtml, tag string) (string, string, error) {
	message := mg.NewMessage(sender, subject, bodyText, recipient)
	message.SetReplyTo(sender)
	err := message.AddTag("Dev-Photo-Contest-" + tag)
	if err != nil {
		return "", "", err
	}
	if bodyHtml != "" {
		message.SetHtml(bodyHtml)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Send the message with a 20 second timeout
	// message, id, err
	return mg.Send(ctx, message)

}

func SendResetEmail(mg *mailgun.MailgunImpl, email string, resetID string) (string, string, error) {
	tag := "Password-Reset"
	bodyText, err := ReadFile("email_templates/password_reset.txt")
	if err != nil {
		return "", "", err
	}
	bodyText = strings.ReplaceAll(bodyText, "<reset_code>", resetID)
	bodyHtml := ""
	return SendEmail(mg, "dnalc-it@cshl.edu", "[DNALC Photo Contest] Password Reset Requested", email, bodyText, bodyHtml, tag)
}
