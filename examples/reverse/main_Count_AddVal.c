/* r2dec pseudo code output */
/* ./reverse @ 0x49aea0 */
#include <stdint.h>

int64_t main_Count_AddVal (int64_t arg_8h, int64_t arg_10h, int64_t arg_18h, int64_t arg_20h, int64_t arg1, int64_t arg4) {
    rdi = arg1;
    rcx = arg4;
    do {
        if (rsp > *((r14 + 0x10))) {
            rdx = *((rbx + 0x30));
            rbx = sym_type_main_CtxCount;
            rax = rcx;
            rcx = obj_runtime_ebss;
            rax = void (*rdx)(uint64_t) (rbx);
            rdx = sym_type__main_Count;
            if (rax == rdx) {
                rax = sym_type_string;
                rbx = 0x004e2698;
                runtime_gopanic ();
            }
            rcx = *((rsp + 0x28));
            *((rcx + 8))++;
            rax = *((rsp + 0x40));
            rax += *((rcx + 0x10));
            *((rcx + 0x10)) = rax;
            ebx = 0xa;
            rax = strconv_FormatInt ();
            ecx = 0;
            edi = 0;
            return rax;
        }
        *((rsp + 8)) = rax;
        *((rsp + 0x10)) = rbx;
        *((rsp + 0x18)) = rcx;
        *((rsp + 0x20)) = rdi;
        runtime_morestack_noctxt_abi0 ();
        rax = *((rsp + 8));
        rbx = *((rsp + 0x10));
        rcx = *((rsp + 0x18));
        rdi = *((rsp + 0x20));
        main_Count_AddVal ();
    } while (1);
}