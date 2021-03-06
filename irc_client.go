package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
)

type IrcConfig struct {
	server string
	port   int
	ssl    bool
}

type IrcClient struct {
	cli      *WSClient
	config   *IrcConfig
	conn     *net.Conn
	conn_tls *tls.Conn

	io *bufio.ReadWriter
}

func (i *IrcClient) Run() {
	fmt.Println("Starting new irc connection to:", i.config.server)
	if i.config.ssl {
		tc := &tls.Config{InsecureSkipVerify: true}
		nc, er := tls.Dial("tcp", fmt.Sprintf("%v:%v", i.config.server, i.config.port), tc)
		if er == nil {
			i.conn_tls = nc

			stid := int32(3)
			stc := int32(1)
			i.cli.SendStatusMessage(stid, stc, i.config.server)

			i.io = bufio.NewReadWriter(bufio.NewReader(i.conn_tls), bufio.NewWriter(i.conn_tls))
			go func() {
				defer i.Close()
				cmdid := int32(2)
				for {
					s, rer := i.io.ReadString('\n')
					if rer == nil {
						fwm := &Command{
							Id: &cmdid,
							ServerMessage: &ServerMessage{
								Server: &i.config.server,
								Msg:    &s,
							},
						}
						i.cli.SendMessage(fwm)
					} else {
						fmt.Println("Read error:", rer.Error())

						stc := int32(2)
						i.cli.SendStatusMessage(stid, stc, i.config.server)
						break
					}
				}
			}()
		} else {
			fmt.Println("Connect error:", er.Error())
		}
	}
}

func (i *IrcClient) SendMessage(msg string) {
	if i.config.ssl {
		_, wer := i.conn_tls.Write([]byte(msg))
		if wer == nil {

		} else {

		}
	}
}

func (i *IrcClient) Close() {
	defer i.cli.RemoveClient(i)
	if i.config.ssl {
		i.conn_tls.Close()
	} else {
		(*i.conn).Close()
	}
}

func NewIrcClient(cli *WSClient, cfg *IrcConfig) *IrcClient {
	return &IrcClient{
		cli:    cli,
		config: cfg,
	}
}
