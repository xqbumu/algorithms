/* r2dec pseudo code output */
/* ./reverse @ 0x49af60 */
#include <stdint.h>

uint64_t main_main (void) {
    int64_t var_28h;
    int64_t var_30h;
    int64_t var_38h;
    int64_t var_40h;
    int64_t var_48h;
    int64_t var_50h;
    int64_t var_58h;
    int64_t var_60h;
    int64_t var_68h;
    int64_t var_70h;
    int64_t var_78h;
    int64_t var_88h;
    int64_t var_90h;
    int64_t var_98h;
    int64_t var_a0h;
    do {
        r12 = rsp - 0x30;
        if (r12 > *((r14 + 0x10))) {
            rax = time_Now ();
            __asm ("movups xmmword [var_58h], xmm15");
            __asm ("movups xmmword [var_68h], xmm15");
            __asm ("movups xmmword [var_78h], xmm15");
            var_60h = rax;
            var_68h = rbx;
            var_70h = rcx;
            rax = &var_58h;
            rbx = obj_go:itab_context_todoCtx_context_Context;
            rcx = obj_runtime_ebss;
            edi = 3;
            main_Count_AddVal ();
            if (rcx == 0) {
                if (rbx != 0) {
                    rax = time_Now ();
                    if (((rax >> 0x3f) & 1) < 0) {
                        rdx = rax;
                        rax <<= 1;
                        rax >>= 0x1f;
                        rsi = 0xdd7b17f80;
                        rbx = rsi + rax;
                    } else {
                        rdx = rax;
                    }
                    rsi = rbx * 0x3b9aca00;
                    edx &= 0x3fffffff;
                    rdx = (int64_t) edx;
                    rdx += rsi;
                    rdx += var_78h;
                    rsi = 0xa1b203eb3d1a0000;
                    rdi = rdx + rsi;
                    rax = &var_58h;
                    rbx = obj_go:itab_context_todoCtx_context_Context;
                    rcx = obj_runtime_ebss;
                    main_Count_AddVal ();
                    if (rcx == 0) {
                        if (rbx != 0) {
                            rax = time_Now ();
                            __asm ("movups xmmword [var_28h], xmm15");
                            __asm ("movups xmmword [var_38h], xmm15");
                            __asm ("movups xmmword [var_48h], xmm15");
                            var_30h = rax;
                            var_38h = rbx;
                            var_40h = rcx;
                            rax = &var_28h;
                            rbx = obj_go:itab_context_todoCtx_context_Context;
                            rcx = obj_runtime_ebss;
                            edi = 3;
                            main_Count_AddVal ();
                            if (rcx != 0) {
                                if (rcx == 0) {
                                } else {
                                    __asm ("movups xmmword [var_88h], xmm15");
                                    __asm ("movups xmmword [var_98h], xmm15");
                                    eax = var_50h;
                                    rax = runtime_convT32 ();
                                    rcx = sym_type_int32;
                                    var_88h = rcx;
                                    var_90h = rax;
                                    rax = var_48h;
                                    rax = runtime_convT64 ();
                                    rcx = sym_type_int64;
                                    var_98h = rcx;
                                    var_a0h = rax;
                                    rax = 0x004be7ef;
                                    ebx = 0x12;
                                    edi = 2;
                                    rsi = rdi;
                                    rcx = &var_88h;
                                    fmt_Sprintf ();
                                    if (rbx != 0) {
                                        rax = time_Now ();
                                        if (((rax >> 0x3f) & 1) < 0) {
                                            rdx = rax;
                                            rax <<= 1;
                                            rax >>= 0x1f;
                                            rsi = 0xdd7b17f80;
                                            rbx = rax + rsi;
                                        } else {
                                            rdx = rax;
                                        }
                                        rsi = rbx * 0x3b9aca00;
                                        edx &= 0x3fffffff;
                                        rdx = (int64_t) edx;
                                        rdx += rsi;
                                        rsi = 0xa1b203eb3d1a0000;
                                        rdi = rsi + rdx;
                                        rax = &var_28h;
                                        rbx = obj_go:itab_context_todoCtx_context_Context;
                                        rcx = obj_runtime_ebss;
                                        rax = main_Count_AddVal ();
                                        if (rcx != 0) {
                                            if (rcx == 0) {
                                            } else {
                                                __asm ("movups xmmword [var_88h], xmm15");
                                                __asm ("movups xmmword [var_98h], xmm15");
                                                eax = var_50h;
                                                rax = runtime_convT32 ();
                                                rcx = sym_type_int32;
                                                var_88h = rcx;
                                                var_90h = rax;
                                                rax = var_48h;
                                                rax = runtime_convT64 ();
                                                rcx = sym_type_int64;
                                                var_98h = rcx;
                                                var_a0h = rax;
                                                rax = 0x004be7ef;
                                                ebx = 0x12;
                                                rcx = &var_88h;
                                                edi = 2;
                                                rsi = rdi;
                                                fmt_Sprintf ();
                                                if (rbx != 0) {
                                                    return rax;
                                                }
                                                rax = sym_type_string;
                                                rbx = 0x004e2708;
                                                runtime_gopanic ();
                                            }
                                            rcx = *((rcx + 8));
                                        }
                                        rax = *((rcx + 8));
                                        rbx = rdi;
                                        runtime_gopanic ();
                                    }
                                    rax = sym_type_string;
                                    rbx = 0x004e2708;
                                    runtime_gopanic ();
                                }
                                rcx = *((rcx + 8));
                            }
                            rax = *((rcx + 8));
                            rbx = rdi;
                            runtime_gopanic ();
                        }
                        rax = sym_type_string;
                        rbx = 0x004e2718;
                        runtime_gopanic ();
                    }
                    if (rbx != 0) {
                        rcx = *((rcx + 8));
                    }
                    rax = *((rcx + 8));
                    rbx = rdi;
                    runtime_gopanic ();
                }
                rax = sym_type_string;
                rbx = 0x004e2718;
                runtime_gopanic ();
            }
            if (rbx != 0) {
                rcx = *((rcx + 8));
            }
            rax = *((rcx + 8));
            rbx = rdi;
            runtime_gopanic ();
        }
        runtime_morestack_noctxt_abi0 ();
        main_main ();
    } while (1);
}