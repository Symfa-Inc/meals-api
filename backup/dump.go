package backup

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Aiscom-LLC/meals-api/src/config"
)

func CreateBackup() {
	fmt.Println("qewqs")
	cmd := exec.Command("pg_dump", "-u", config.Env.DbUser, "-w", config.Env.DbName)
	cmd.Stdin = strings.NewReader(config.Env.DbPassword)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err == nil {
		file, err := os.Create("dumps/dump-" + time.Now().UTC().String() + ".dump")
		if err != nil {
			log.Println(err)
		}

		_, err = file.WriteString(out.String())
		if err != nil {
			log.Println(err)
		}
	}
}
