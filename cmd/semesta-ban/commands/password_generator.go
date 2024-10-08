package commands

import (
	"fmt"
	"libra-internal/bootstrap"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	registerCommand(generatePassword)
}

// todo create script for  migration db automation
func generatePassword(dep *bootstrap.Dependency) *cobra.Command {
	return &cobra.Command{
		Use:   "password",
		Short: "Generate password",
		Long:  `This command is used to start generate password service`,
		Run: func(cmd *cobra.Command, args []string) {

			// start := time.Now()
			// a := helper.GenerateTransactionId("", start.Format("20060102"))
			// fmt.Println(a)
			// log.Info("Server shutdown gracefully.")
			// // fmt.Println(start.Format("2006-01-02"))

			// end := start.AddDate(0, 0, 30)
			// // fmt.Println(end.Format("2006-01-02"))
			// i := 1
			// for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			// 	fmt.Println(i)
			// 	fmt.Println(d.Format("2006-01-02"))
			// 	i++
			// }

			// a := "asd hehe KOK DQW-sad asdsad DDDASD aaDD diprosesS"
			// a = strings.Title(strings.ToLower(a))
			// fmt.Println(a)

			// p := message.NewPrinter(language.English)
			// test := p.Sprintf("%d\n", 2000000)
			// fix := strings.ReplaceAll(test, ",", ".")
			// fmt.Println(fix)

			// sEnc := helper.GenerateB64AuthMidtrans("SB-Mid-server-cXRk9vIv_uoZ0nfWgHnqozI7")
			// fmt.Println(sEnc)

			hashedPass, err := bcrypt.GenerateFromPassword([]byte("sugiono"), bcrypt.DefaultCost)
			if err != nil {
				return
			}
			fmt.Println(string(hashedPass))

			// e := os.Remove("suzuki.png")
			// if e != nil {
			// 	log.Fatal(e)
			// }

		},
	}
}
