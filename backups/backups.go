package backups

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Aiscom-LLC/meals-api/config"
)

func CreateBackup() {
	var out bytes.Buffer

	cmd := exec.Command("pg_dump", "-u", config.Env.DbUser, "-w", config.Env.DbName)

	cmd.Stdin = strings.NewReader(config.Env.DbPassword)
	cmd.Stdout = &out

	err := cmd.Run()

	if err == nil {
		file, err := os.Create("backups/dumps/dump-" + time.Now().UTC().String() + ".psql")

		if err != nil {
			log.Println(err)
		}

		_, err = file.WriteString(out.String())
		if err != nil {
			log.Println(err)
		}
	}
}
