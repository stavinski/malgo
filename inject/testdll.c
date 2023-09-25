/*

Just simple test DLL code to test the injection with.

compile with:

gcc -c -o testdll.o .\testdll.c -D ADD_EXPORTS
gcc -o testdll.dll testdll.o -s -shared -lkernel32

use:

inject.exe c:\<path>\testdll.dll <pid>

*/

#define WIN32_LEAN_AND_MEAN
#define _WIN32_WINNT 0x0500

#if _MSC_VER
#pragma comment(lib,"user32.lib")
#endif
#pragma comment(lib,"kernel32.lib")

#include <Windows.h>

// Main entry point
BOOL WINAPI DllMain(HINSTANCE hinstDLL, DWORD ul_reason_for_call, LPVOID lpvReserved)
{
    if (ul_reason_for_call == DLL_PROCESS_ATTACH)
    {
        MessageBoxA(0, "Hello from test DLL", "Hello", 0);
        return TRUE;
    }
}