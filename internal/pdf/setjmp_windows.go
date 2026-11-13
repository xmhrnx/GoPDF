package pdf

// The MuPDF library go-fitz ships for Windows was compiled against the UCRT
// C runtime, where setjmp expands to __intrinsic_setjmpex, a function only
// ucrtbase.dll exports. MinGW-w64 toolchains that link the older msvcrt
// runtime (MSYS2's MINGW64, most distributions' cross compilers) don't
// define it, and the link fails with "undefined reference to
// __intrinsic_setjmpex".
//
// msvcrt.dll exports the same function as _setjmpex: same arguments, same
// jmp_buf layout, paired with the same longjmp. On those toolchains this
// defines __intrinsic_setjmpex as a tail jump to it. A jump, not a call,
// keeps the caller's stack frame intact, which setjmp must record. UCRT
// toolchains (where _mingw.h defines _UCRT) already have the real one.

/*
#include <_mingw.h>

#if defined(__x86_64__) && !defined(_UCRT)
__asm__(
	".text\n"
	".globl __intrinsic_setjmpex\n"
	".def __intrinsic_setjmpex; .scl 2; .type 32; .endef\n"
	"__intrinsic_setjmpex:\n"
	"	jmp *__imp__setjmpex(%rip)\n"
);
#endif
*/
import "C"
