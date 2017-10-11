package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"github.com/robfig/cron"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug = kingpin.Flag("debug", "").Short('d').Bool()
	token = kingpin.Flag("token", "").Short('t').Envar("SLACK_TOKEN").String()
)

var (
	users             = make(map[string]map[string][]string)
	latestReleasesIDs = make(map[string]map[string]int)
	channels          = make(map[string]string)
)

func main() {
	kingpin.Version("0.1.0")

	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.VersionFlag.Short('v')

	kingpin.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	cron := cron.New()
	cron.Start()
	defer cron.Stop()

	api := slack.New(*token)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	client := github.NewClient(nil)
	ctx := context.Background()

	cron.AddFunc("*/30 * * * * ?", func() {
		log.Debug("tick")

		for owner, _ := range users {
			for repo, _ := range users[owner] {
				rel, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
				log.WithFields(log.Fields{
					"owner": owner,
					"repo":  repo,
					"ID":    *rel.ID,
				}).Debug("Latest release ID")

				if err != nil {
					log.Error(err)
					break
				}

				if latestReleasesIDs[owner][repo] != 0 && latestReleasesIDs[owner][repo] != *rel.ID {
					log.WithFields(log.Fields{
						"owner":  owner,
						"repo":   repo,
						"prevID": latestReleasesIDs[owner][repo],
						"nextID": *rel.ID,
					}).Debug("Replacing latest release ID")

					latestReleasesIDs[owner][repo] = *rel.ID
					for _, user := range users[owner][repo] {
						rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("hey! new release is available %s!", *rel.HTMLURL), channels[user]))
					}
				}
			}
		}
	})

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			log.WithFields(log.Fields{
				"User":    ev.User,
				"Text":    ev.Text,
				"Channel": ev.Channel,
			}).Debug("Message")

			user := ev.User
			channels[user] = ev.Channel

			s := regexp.MustCompile("^((?i)[a-z]+).*(?:[[:blank:]]+|github.com[/:])([[:alnum:]][[:alnum:]-]*[[:alnum:]])/([[:alnum:]-_]+)").FindStringSubmatch(ev.Text)
			if s == nil {
				log.Debug(s)
				break
			}
			cmd, owner, repo := strings.ToLower(s[1]), s[2], s[3]

		Cmd:
			switch cmd {
			case "watch":
				log.WithFields(log.Fields{
					"user":  user,
					"owner": owner,
					"repo":  repo,
				}).Debug("Watch")

				if users[owner] == nil {
					users[owner] = make(map[string][]string)
				}

				for _, v := range users[owner][repo] {
					if v == user {
						break Cmd
					}
				}

				users[owner][repo] = append(users[owner][repo], user)
				rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("watching https://github.com/%s/%s for you!", owner, repo), ev.Channel))

				if latestReleasesIDs[owner] == nil {
					latestReleasesIDs[owner] = make(map[string]int)
				}

				rel, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
				if err != nil {
					log.Error(err)
					break Cmd
				}

				log.WithFields(log.Fields{
					"owner": owner,
					"repo":  repo,
					"ID":    *rel.ID,
				}).Debug("Latest release ID")
				latestReleasesIDs[owner][repo] = *rel.ID

			case "unwatch":
				log.WithFields(log.Fields{
					"user":  user,
					"owner": owner,
					"repo":  repo,
				}).Debug("Unwatch")

				if users[owner] == nil {
					break Cmd
				}

				for i, v := range users[owner][repo] {
					if v == user {
						users[owner][repo] = append(users[owner][repo][:i], users[owner][repo][i+1:]...)
						rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("unwatched https://github.com/%s/%s", owner, repo), ev.Channel))

						if len(users[owner][repo]) == 0 {
							delete(users[owner], repo)
						}

						if len(users[owner]) == 0 {
							delete(users, owner)
						}
					}
				}
			}

		case *slack.RTMError:
			log.Error(ev.Error())

		case *slack.InvalidAuthEvent:
			log.Error("Invalid authentication token.")
			return

		default:
		}
	}
}
