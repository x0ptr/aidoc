package format

import "fmt"

func PrintBanner(title string) string {
	return fmt.Sprintf("\n\033[1;37;46m doc of: %s  \033[0m\n\n", title)
}
