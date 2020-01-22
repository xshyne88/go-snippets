package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"syscall"
	"time"
)

type LeaderBoard struct {
	scores       map[int]int
	broadcastOut chan string
}

func (l *LeaderBoard) startBroadcastStream() {
	for {
		select {
		case m := <-l.broadcastOut:
			fmt.Print(m)
		default:
		}
	}
}

func newLeaderBoard() *LeaderBoard {
	l := &LeaderBoard{make(map[int]int), make(chan string)}
	go l.startBroadcastStream()
	return l
}

func (l *LeaderBoard) Scores() map[int]int {
	l.broadcastOut <- fmt.Sprint("Reporting Scores: \n")
	return l.scores
}

func (l *LeaderBoard) AdjustScore(playerID int, newScore int) {
	l.broadcastOut <- fmt.Sprintf("Player %d score changed to %d: \n", playerID, newScore)
}

func (l *LeaderBoard) AddPlayer(playerID int) {
	l.scores[playerID] = 0
	l.broadcastOut <- fmt.Sprintf("Player %d as added \n", playerID)
}

func (l *LeaderBoard) DeletePlayer(playerID int) {
	delete(l.scores, playerID)
	l.broadcastOut <- fmt.Sprintf("Player %d was removed \n", playerID)
}

func main() {
	leaderBoardChan := make(chan *LeaderBoard)
	l := newLeaderBoard()
	go func() {
		leaderBoardChan <- l
	}()

	go addPlayer2(leaderBoardChan)
	go sleepAddPlayer(1, leaderBoardChan)
	go sleepAddPlayer(2, leaderBoardChan)
	go sleepAdjustScore(3, 2, 50, leaderBoardChan)
	go sleepDeletePlayer(4, 2, leaderBoardChan)
	go sleepAddPlayer(4, leaderBoardChan)
	go sleepAddPlayer(5, leaderBoardChan)
	go sleepShowScores(7, leaderBoardChan)
	go killAfterEightSeconds()

	signals := make(chan os.Signal, 1)
	go func() {
		for {
			select {
			case sig, _ := <-signals:
				switch sig {
				case syscall.SIGTERM:
					log.Printf("CLOSING GRACEFULLY **")
					os.Exit(0)
				}
			default:
			}
		}
	}()

	select {}
}

func addPlayer2(ch chan *LeaderBoard) {
	l := <-ch
	l.AddPlayer(2)
	ch <- l
}

func killAfterEightSeconds() {
	sleep(8)
	os.Exit(0)
}

func sleep(duration int) {
	time.Sleep(time.Second * time.Duration(duration))
}

func sleepDeletePlayer(duration int, playerID int, ch chan *LeaderBoard) {
	sleep(duration)
	lb := <-ch
	lb.DeletePlayer(playerID)
	ch <- lb
}

func sleepAddPlayer(duration int, ch chan *LeaderBoard) {
	sleep(duration)
	lb := <-ch
	lb.AddPlayer(rand.Intn(200))
	ch <- lb
}

func sleepAdjustScore(duration, playerID, score int, ch chan *LeaderBoard) {
	sleep(duration)
	lb := <-ch
	lb.AdjustScore(playerID, score)
	ch <- lb
}

func sleepShowScores(duration int, ch chan *LeaderBoard) {
	sleep(duration)
	lb := <-ch
	fmt.Print(lb.Scores())
	ch <- lb
}
