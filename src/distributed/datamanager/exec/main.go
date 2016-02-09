package exec
import (
	"distributed/qutils"
	"log"
	"bytes"
	"encoding/gob"
	"distributed/dto"
	"distributed/datamanager"
)

const url = "amqp://guest:guest@localhost:5672"

func main() {
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		qutils.PersistReadingsQueue,
		"",
		false,
		true,
		false,
		false,
		nil)

	if err != nil {
		log.Fatalln("Failed to get access to messages")
	}

	for msg := range msgs {
		buf := bytes.NewBuffer(msg.Body)
		dec := gob.NewDecoder(buf)
		sd := &dto.SensorMessage{}
		dec.Decode(sd)

		err := datamanager.SaveReadings(sd)

		if err != nil {
			log.Printf("Failed to save reading from sensor %v. Error: %s",
				sd.Name, err.Error())
		} else {
			msg.Ack(false)
		}
	}
}