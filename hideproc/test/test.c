#include <stdlib.h>
#include <stdio.h>
#include <Windows.h>
#include <psapi.h>
#include <winternl.h>


/*

Very simple C windows app to test IAT hook against. It runs the same code twice to enumerate over processes, the msgbox gives a chance to then inject the hideproc.dll DLL and see that the explorer.exe process(es)
are now not shown in the results.

compile with:

gcc -o test.exe test.c -lntdll

*/

typedef struct _SYSTEM_PROCESS_INFO
{
    ULONG                   NextEntryOffset;
    ULONG                   NumberOfThreads;
    LARGE_INTEGER           Reserved[3];
    LARGE_INTEGER           CreateTime;
    LARGE_INTEGER           UserTime;
    LARGE_INTEGER           KernelTime;
    UNICODE_STRING          ImageName;
    ULONG                   BasePriority;
    HANDLE                  ProcessId;
    HANDLE                  InheritedFromProcessId;
} SYSTEM_PROCESS_INFO,*PSYSTEM_PROCESS_INFO;


BOOL PrintProcesses()
{
    PSYSTEM_PROCESS_INFO procInfo = (PSYSTEM_PROCESS_INFO)VirtualAlloc(NULL, 1024*1024, MEM_COMMIT|MEM_RESERVE,PAGE_READWRITE);
    ULONG sysinfoLength = 1024*1024;
    ULONG returnLength;
    PWORD imageName;

    if (!procInfo)
    {
        return FALSE;
    }

    NTSTATUS status = NtQuerySystemInformation(SystemProcessInformation, procInfo, sysinfoLength, &returnLength);
    if (!NT_SUCCESS(status))
    {
        VirtualFree(procInfo, 0, MEM_RELEASE);
        return FALSE;
    }

    while(procInfo->NextEntryOffset != 0) // Loop over the list until we reach the last entry.
    {
        PWSTR imageName = (procInfo->ProcessId ==0) ? L"Idle" : procInfo->ImageName.Buffer;
        char *lpConverted = "";
        int bufferSize = WideCharToMultiByte(CP_UTF8, WC_DEFAULTCHAR, imageName, -1, lpConverted, 0, NULL, NULL);
        lpConverted = malloc(bufferSize);
        if (lpConverted != NULL){
            WideCharToMultiByte(CP_UTF8, WC_DEFAULTCHAR, imageName, -1, lpConverted, bufferSize, NULL, NULL);
            if (strcmpi(lpConverted, "explorer.exe") == 0){
                printf("Process name: %s | Process ID: %d\n", lpConverted,procInfo->ProcessId); // Display process information.
            }
            // printf("Process name: %s | Process ID: %d\n", lpConverted,procInfo->ProcessId); // Display process information.
            free(lpConverted);
        }

        procInfo=(PSYSTEM_PROCESS_INFO)((LPBYTE)procInfo + procInfo->NextEntryOffset); // Calculate the address of the next entry.
    }

    VirtualFree(procInfo, 0, MEM_RELEASE);
    return TRUE;
}

// Entry point
int main(){
    if (!PrintProcesses())
    {
        fprintf(stderr, "Could not get processes.\n");
        return 1;
    }

    MessageBoxA(0, "Install hook now...", "TEST", 0);
    printf("Processes after hook...\n");

    if (!PrintProcesses())
    {
        fprintf(stderr, "Could not get processes.\n");
        return 1;
    }

    return 0;
}