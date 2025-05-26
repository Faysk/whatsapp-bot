package utils

import (
	"log"
	"os"
)

// SetupLogger configura o logger padrão (por enquanto simples, mas extensível)
func SetupLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("📝 Logger inicializado.")
}
