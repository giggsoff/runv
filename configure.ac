#                                               -*- Autoconf -*-
# Process this file with autoconf to produce a configure script.

AC_PREREQ([2.69])
AC_INIT([runv], [1.0.0], [hyper.sh])
AM_INIT_AUTOMAKE([-Wall -Werror foreign subdir-objects])

# Checks for programs.
AC_PROG_CC
AM_PROG_AR
AC_PROG_RANLIB
AC_CANONICAL_BUILD

# Checks for go tool chain
AC_CHECK_PROG([has_go], [go], [yes], [no])
if test "x$has_go" != "xyes" ; then
    AC_MSG_ERROR(Unable to find go binary in PATH)
fi

# Platform specific setup
AC_CANONICAL_HOST

case $host_os in
	linux*) AM_CONDITIONAL([ON_LINUX], [ true  ]) ;;
	*)      AM_CONDITIONAL([ON_LINUX], [ false ]) ;;
esac

# Check for which host we are on and setup a few things
# specifically based on the host


# Checks for libraries.

LIBVIRT_REQUIRED="1.2.2"

AC_ARG_WITH([libvirt],
            [AS_HELP_STRING([--without-libvirt],
                            [run runv with libvirt])],
            [with_libvirt=no],[with_libvirt=yes])

if test "$with_libvirt" = yes; then
	PKG_CHECK_MODULES([LIBVIRT], [libvirt >= $LIBVIRT_REQUIRED], [], [with_libvirt=no])
fi

AM_CONDITIONAL([WITH_LIBVIRT], [test "x$with_libvirt" == "xyes"])

# Checks for header files.
AC_CHECK_HEADERS([stdlib.h string.h libxl.h])

# Checks for typedefs, structures, and compiler characteristics.
AC_CHECK_HEADER_STDBOOL
AC_TYPE_PID_T
AC_TYPE_UINT32_T
AC_TYPE_UINT8_T

# Checks for library functions.
AC_FUNC_MALLOC
AC_CHECK_FUNCS([strdup])

AC_ARG_WITH([xen],
            [AS_HELP_STRING([--without-xen],
                            [run runv with xen (libxl, need xen 4.5 or higher)])],
            [with_xen=no],[with_xen=yes])

if test "x$with_xen" != "xno" ; then
    # Checks for libxl
    AC_CHECK_HEADERS([libxl.h libxl_utils.h], [libxl_found=yes], [libxl_found=no])
    if test "x$libxl_found" != "xyes"; then
           with_xen=no
    fi
fi

XEN_REQUIRED="4.5.0"
XENPV_REQUIRED="4.9.0"

if test "$with_xen" = yes; then
    PKG_CHECK_MODULES([xen_4_5], [xenlight >= $XEN_REQUIRED], [with_xen=yes], [with_xen=no])
    PKG_CHECK_MODULES([xen_4_9], [xenlight >= $XENPV_REQUIRED], [with_xen490=yes], [with_xen490=no])
fi

if test "x$with_xen" != "xno" ; then
    AC_DEFINE_UNQUOTED([WITH_XEN], 1, [run runv with xen])
fi
if test "x$with_xen490" = "xyes" ; then
    AC_DEFINE_UNQUOTED([WITH_XEN490], 1, [run runv with xen 4.9 or above])
fi

AM_CONDITIONAL([WITH_XEN], [test "x$with_xen" == "xyes"])
AM_CONDITIONAL([WITH_XEN490], [test "x$with_xen490" == "xyes"])

AC_CONFIG_FILES([Makefile])

AC_OUTPUT

AC_MSG_RESULT([
	${PACKAGE}  ${VERSION}
	build OS:   ${build_os}
	prefix:	    ${prefix}

	has go:     ${has_go}

	with xen:   ${with_xen}

	with libvirt: ${with_libvirt}
])
