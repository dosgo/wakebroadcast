// +build windows
package notify

import (

	"github.com/lxn/walk"
	_"github.com/lxn/walk"
	"log"
	"os"
)
var mw *walk.MainWindow

func init(){
	var err error
	mw,  err= walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
}

func MsgBox(title string,msg string ){
	walk.MsgBox(mw, title, msg, walk.MsgBoxIconInformation)
}

func GuiInit()  {
	//托盘图标文件
	icon, err := walk.Resources.Icon("./icon.ico")
	if err != nil {
		log.Fatal(err)
	}
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip("唤醒广播转发服务运行中"); err != nil {
		log.Fatal(err)
	}

	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		if err := ni.ShowInfo("唤醒广播转发服务运行中","右键退出"); err != nil {
			log.Fatal(err)
		}
	})
	exitAction := walk.NewAction()
	if err := exitAction.SetText("退出"); err != nil {
		log.Fatal(err)
	}
	//Exit 实现的功能
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}
	if err := ni.ShowInfo("Walk NotifyIcon Example", "Click the icon to show again."); err != nil {
		log.Fatal(err)
	}
	mw.Run()
}


