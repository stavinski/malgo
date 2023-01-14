#include <Windows.h>

/*

Very simple C windows app to test IAT hook against. Run and then when the first msgbox is shown inject the iathook.dll into the process and click OK
on the msgbox, you should see that the call has been hooked and presented with a msgbox from the hook code.

compile with:

gcc -o test.exe test.c

*/
int main(){
    MessageBoxA(NULL, "Text before hook, place hook now...", "TEST", 0);
    MessageBoxA(NULL, "Text after hook", "TEST", 0);
    return 0;
}