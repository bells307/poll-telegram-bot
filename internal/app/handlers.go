package app

import (
	"log"

	tm "github.com/and3rson/telemux/v2"
)

func (b *PollBot) HelpHandler(u *tm.Update) {
	log.Println("HelpHandler")
}

func (b *PollBot) PollHandler(u *tm.Update) {
	log.Println("PollHandler")
}

func (b *PollBot) AddPollOptionHandler(u *tm.Update) {
	log.Println("AddPollOptionHandler")
}

func (b *PollBot) DeletePollOptionHandler(u *tm.Update) {
	log.Println("DeletePollOptionHandler")
}

func (b *PollBot) ListPollOptionsHandler(u *tm.Update) {
	log.Println("ListPollOptionsHandler")
}
