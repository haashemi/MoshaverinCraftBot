package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/LlamaNite/llamalog"
	"github.com/haashemi/MoshaverinCraftBot/internal/config"
	"github.com/haashemi/MoshaverinCraftBot/internal/ipapi"
	"github.com/haashemi/tgo"
	"github.com/haashemi/tgo/filters"
	"github.com/haashemi/tgo/routers/message"
)

var log = llamalog.NewLogger("IPBot")

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load the config file > %v", err)
	}

	var transport http.RoundTripper
	if conf.Proxy != "" {
		proxy, err := url.Parse(conf.Proxy)
		if err != nil {
			log.Fatal("Failed to parse proxy URL: %v", err)
		}

		transport = &http.Transport{Proxy: http.ProxyURL(proxy)}

		log.Info("Proxy URL: %v", proxy)
	}

	bot := tgo.NewBot(conf.Token, tgo.Options{
		Client:           &http.Client{Transport: transport},
		DefaultParseMode: tgo.ParseModeHTML,
	})

	info, err := bot.GetMe()
	if err != nil {
		log.Fatal("Failed to get bot info: %v", err)
	}
	log.Info("Bot identified as %s", info.Username)

	mr := message.NewRouter()
	mr.Handle(filters.And(filters.Command("start", info.Username)), handleStart)
	mr.Handle(filters.And(filters.Command("ip", info.Username), filters.Whitelist(conf.Whitelist...)), handleIP)
	_ = bot.AddRouter(mr)

	_, _ = bot.SetMyCommands(&tgo.SetMyCommands{
		Commands: []*tgo.BotCommand{
			{Command: "start", Description: "سلام"},
			{Command: "ip", Description: "دریافت آیپی"},
		},
	})

	for {
		log.Info("Polling started.")
		err = bot.StartPolling(30)
		log.Error("Polling failed > %v", err)
		log.Warn("Bot stopped. Restarting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

func handleStart(ctx *message.Context) {
	_, _ = ctx.Send(&tgo.SendMessage{Text: strings.TrimSpace(`
درود!

این ربات مخصوص مشاورینه برای اینکه از آیپی سرور ماینکرفت مطلع بشن.
کل کاری که لازمه بکنی اینه که کامند /ip رو بفرستی.

اگه دیدی بهت اهمیت نداد به صاحبش پیام بده ببینیم چه خبر شده.
`)})
}

func handleIP(ctx *message.Context) {
	_, _ = ctx.Bot.API.SendChatAction(&tgo.SendChatAction{ChatId: tgo.ID(ctx.Chat.Id), Action: "typing"})

	ip, err := ipapi.GetIP()
	if err != nil {
		ctx.Send(&tgo.SendMessage{Text: "دریافت آیپی با موفقیت با شکست مواجه شد."})
		log.Error("Failed to get IP: %v", err)
		return
	}

	_, _ = ctx.Send(&tgo.SendMessage{Text: "<code>" + ip.Query + "</code>"})
}
