package app

import (
	"github.com/golang/protobuf/proto"
	"github.com/rfw141/i/log"
	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
	"regexp"
	"strings"
)

type ExecHandler func(ctx Ctx) (proto.Message, error)

type Executor struct {
	Path        string
	ExecHandler ExecHandler
}

type Cmd struct {
	*Cfg
	root *cobra.Command
	rpcs map[string]*Rpc
}

func NewCmd(cfg *Cfg, rpcs map[string]*Rpc) *Cmd {
	var c Cmd
	c.Cfg = cfg
	c.root = &cobra.Command{
		Use: cfg.Name,
	}
	c.rpcs = rpcs
	return &c
}

func (c *Cmd) RegisterExecutors(executors []Executor) {
	for _, executor := range executors {
		nExecutor := executor
		use := nExecutor.Path[strings.LastIndex(nExecutor.Path, "/")+1:]
		log.Debugf("(cmd) register %s", use)
		c.root.AddCommand(&cobra.Command{
			Use: use,
			Run: func(cmd *cobra.Command, args []string) {
				if cmd.Flag("example").Value.String() == "true" {
					if rpc, ok := c.rpcs[nExecutor.Path]; ok {
						rpc.Example()
						return
					}
					log.Errorf("exec %s error: not found", use)
					return
				}
				ctx := NewCmdCtx()
				ctx.req = parseArgs(args)
				rsp, err := nExecutor.ExecHandler(ctx)
				if err != nil {
					log.Errorf("exec %s error: %v", use, err)
					return
				}
				log.Infof("exec %s success: %+v", use, rsp)
			},
		})
	}
}

func (c *Cmd) Execute(_ Ctx) error {
	return c.root.Execute()
}

func parseArgs(args []string) string {
	var req string
	for _, v := range args {
		vs := strings.Split(v, "=")
		if len(vs) == 2 {
			key, val := vs[0], vs[1]
			switch {
			case val == "true" || val == "false",
				strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\""),
				strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}"),
				strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]"),
				regexp.MustCompile(`\d`).MatchString(val):
				req, _ = sjson.SetRaw(req, key, val)
			default:
				req, _ = sjson.Set(req, key, val)
			}
		}
	}
	return req
}
