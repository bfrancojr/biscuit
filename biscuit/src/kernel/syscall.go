package main

import "fmt"
import "math/rand"
import "runtime"
import "runtime/debug"
import "runtime/pprof"
import "sort"
import "sync"
import "sync/atomic"
import "time"
import "unsafe"

import "bounds"
import "bpath"
import "common"
import "defs"
import "fs"
import "limits"
import "mem"
import "res"
import "stat"
import "tinfo"
import "ustr"
import "util"
import "vm"

var _sysbounds = []int{
	//var _sysbounds = map[int]int {
	common.SYS_READ:       bounds.Bounds(bounds.B_SYS_READ),
	common.SYS_WRITE:      bounds.Bounds(bounds.B_SYS_WRITE),
	common.SYS_OPEN:       bounds.Bounds(bounds.B_SYS_OPEN),
	common.SYS_CLOSE:      bounds.Bounds(bounds.B_SYSCALL_T_SYS_CLOSE),
	common.SYS_STAT:       bounds.Bounds(bounds.B_SYS_STAT),
	common.SYS_FSTAT:      bounds.Bounds(bounds.B_SYS_FSTAT),
	common.SYS_POLL:       bounds.Bounds(bounds.B_SYS_POLL),
	common.SYS_LSEEK:      bounds.Bounds(bounds.B_SYS_LSEEK),
	common.SYS_MMAP:       bounds.Bounds(bounds.B_SYS_MMAP),
	common.SYS_MUNMAP:     bounds.Bounds(bounds.B_SYS_MUNMAP),
	common.SYS_SIGACT:     bounds.Bounds(bounds.B_SYS_SIGACTION),
	common.SYS_READV:      bounds.Bounds(bounds.B_SYS_READV),
	common.SYS_WRITEV:     bounds.Bounds(bounds.B_SYS_WRITEV),
	common.SYS_ACCESS:     bounds.Bounds(bounds.B_SYS_ACCESS),
	common.SYS_DUP2:       bounds.Bounds(bounds.B_SYS_DUP2),
	common.SYS_PAUSE:      bounds.Bounds(bounds.B_SYS_PAUSE),
	common.SYS_GETPID:     bounds.Bounds(bounds.B_SYS_GETPID),
	common.SYS_GETPPID:    bounds.Bounds(bounds.B_SYS_GETPPID),
	common.SYS_SOCKET:     bounds.Bounds(bounds.B_SYS_SOCKET),
	common.SYS_CONNECT:    bounds.Bounds(bounds.B_SYS_CONNECT),
	common.SYS_ACCEPT:     bounds.Bounds(bounds.B_SYS_ACCEPT),
	common.SYS_SENDTO:     bounds.Bounds(bounds.B_SYS_SENDTO),
	common.SYS_RECVFROM:   bounds.Bounds(bounds.B_SYS_RECVFROM),
	common.SYS_SOCKPAIR:   bounds.Bounds(bounds.B_SYS_SOCKETPAIR),
	common.SYS_SHUTDOWN:   bounds.Bounds(bounds.B_SYS_SHUTDOWN),
	common.SYS_BIND:       bounds.Bounds(bounds.B_SYS_BIND),
	common.SYS_LISTEN:     bounds.Bounds(bounds.B_SYS_LISTEN),
	common.SYS_RECVMSG:    bounds.Bounds(bounds.B_SYS_RECVMSG),
	common.SYS_SENDMSG:    bounds.Bounds(bounds.B_SYS_SENDMSG),
	common.SYS_GETSOCKOPT: bounds.Bounds(bounds.B_SYS_GETSOCKOPT),
	common.SYS_SETSOCKOPT: bounds.Bounds(bounds.B_SYS_SETSOCKOPT),
	common.SYS_FORK:       bounds.Bounds(bounds.B_SYS_FORK),
	common.SYS_EXECV:      bounds.Bounds(bounds.B_SYS_EXECV),
	common.SYS_EXIT:       bounds.Bounds(bounds.B_SYSCALL_T_SYS_EXIT),
	common.SYS_WAIT4:      bounds.Bounds(bounds.B_SYS_WAIT4),
	common.SYS_KILL:       bounds.Bounds(bounds.B_SYS_KILL),
	common.SYS_FCNTL:      bounds.Bounds(bounds.B_SYS_FCNTL),
	common.SYS_TRUNC:      bounds.Bounds(bounds.B_SYS_TRUNCATE),
	common.SYS_FTRUNC:     bounds.Bounds(bounds.B_SYS_FTRUNCATE),
	common.SYS_GETCWD:     bounds.Bounds(bounds.B_SYS_GETCWD),
	common.SYS_CHDIR:      bounds.Bounds(bounds.B_SYS_CHDIR),
	common.SYS_RENAME:     bounds.Bounds(bounds.B_SYS_RENAME),
	common.SYS_MKDIR:      bounds.Bounds(bounds.B_SYS_MKDIR),
	common.SYS_LINK:       bounds.Bounds(bounds.B_SYS_LINK),
	common.SYS_UNLINK:     bounds.Bounds(bounds.B_SYS_UNLINK),
	common.SYS_GETTOD:     bounds.Bounds(bounds.B_SYS_GETTIMEOFDAY),
	common.SYS_GETRLMT:    bounds.Bounds(bounds.B_SYS_GETRLIMIT),
	common.SYS_GETRUSG:    bounds.Bounds(bounds.B_SYS_GETRUSAGE),
	common.SYS_MKNOD:      bounds.Bounds(bounds.B_SYS_MKNOD),
	common.SYS_SETRLMT:    bounds.Bounds(bounds.B_SYS_SETRLIMIT),
	common.SYS_SYNC:       bounds.Bounds(bounds.B_SYS_SYNC),
	common.SYS_REBOOT:     bounds.Bounds(bounds.B_SYS_REBOOT),
	common.SYS_NANOSLEEP:  bounds.Bounds(bounds.B_SYS_NANOSLEEP),
	common.SYS_PIPE2:      bounds.Bounds(bounds.B_SYS_PIPE2),
	common.SYS_PROF:       bounds.Bounds(bounds.B_SYS_PROF),
	common.SYS_THREXIT:    bounds.Bounds(bounds.B_SYS_THREXIT),
	common.SYS_INFO:       bounds.Bounds(bounds.B_SYS_INFO),
	common.SYS_PREAD:      bounds.Bounds(bounds.B_SYS_PREAD),
	common.SYS_PWRITE:     bounds.Bounds(bounds.B_SYS_PWRITE),
	common.SYS_FUTEX:      bounds.Bounds(bounds.B_SYS_FUTEX),
	common.SYS_GETTID:     bounds.Bounds(bounds.B_SYS_GETTID),
}

// Implements Syscall_i
type syscall_t struct {
}

var sys = &syscall_t{}

func (s *syscall_t) Syscall(p *common.Proc_t, tid defs.Tid_t, tf *[common.TFSIZE]uintptr) int {

	if p.Doomed() {
		// this process has been killed
		p.Reap_doomed(tid)
		return 0
	}

	sysno := int(tf[common.TF_RAX])

	//lim, ok := _sysbounds[sysno]
	//if !ok {
	//	panic("bad limit")
	//}
	lim := _sysbounds[sysno]
	//if lim == 0 {
	//	panic("bad limit")
	//}
	if !res.Resadd(lim) {
		//fmt.Printf("syscall res failed\n")
		return int(-defs.ENOHEAP)
	}

	a1 := int(tf[common.TF_RDI])
	a2 := int(tf[common.TF_RSI])
	a3 := int(tf[common.TF_RDX])
	a4 := int(tf[common.TF_RCX])
	a5 := int(tf[common.TF_R8])

	var ret int
	switch sysno {
	case common.SYS_READ:
		ret = sys_read(p, a1, a2, a3)
	case common.SYS_WRITE:
		ret = sys_write(p, a1, a2, a3)
	case common.SYS_OPEN:
		ret = sys_open(p, a1, a2, a3)
	case common.SYS_CLOSE:
		ret = s.Sys_close(p, a1)
	case common.SYS_STAT:
		ret = sys_stat(p, a1, a2)
	case common.SYS_FSTAT:
		ret = sys_fstat(p, a1, a2)
	case common.SYS_POLL:
		ret = sys_poll(p, tid, a1, a2, a3)
	case common.SYS_LSEEK:
		ret = sys_lseek(p, a1, a2, a3)
	case common.SYS_MMAP:
		ret = sys_mmap(p, a1, a2, a3, a4, a5)
	case common.SYS_MUNMAP:
		ret = sys_munmap(p, a1, a2)
	case common.SYS_READV:
		ret = sys_readv(p, a1, a2, a3)
	case common.SYS_WRITEV:
		ret = sys_writev(p, a1, a2, a3)
	case common.SYS_SIGACT:
		ret = sys_sigaction(p, a1, a2, a3)
	case common.SYS_ACCESS:
		ret = sys_access(p, a1, a2)
	case common.SYS_DUP2:
		ret = sys_dup2(p, a1, a2)
	case common.SYS_PAUSE:
		ret = sys_pause(p)
	case common.SYS_GETPID:
		ret = sys_getpid(p, tid)
	case common.SYS_GETPPID:
		ret = sys_getppid(p, tid)
	case common.SYS_SOCKET:
		ret = sys_socket(p, a1, a2, a3)
	case common.SYS_CONNECT:
		ret = sys_connect(p, a1, a2, a3)
	case common.SYS_ACCEPT:
		ret = sys_accept(p, a1, a2, a3)
	case common.SYS_SENDTO:
		ret = sys_sendto(p, a1, a2, a3, a4, a5)
	case common.SYS_RECVFROM:
		ret = sys_recvfrom(p, a1, a2, a3, a4, a5)
	case common.SYS_SOCKPAIR:
		ret = sys_socketpair(p, a1, a2, a3, a4)
	case common.SYS_SHUTDOWN:
		ret = sys_shutdown(p, a1, a2)
	case common.SYS_BIND:
		ret = sys_bind(p, a1, a2, a3)
	case common.SYS_LISTEN:
		ret = sys_listen(p, a1, a2)
	case common.SYS_RECVMSG:
		ret = sys_recvmsg(p, a1, a2, a3)
	case common.SYS_SENDMSG:
		ret = sys_sendmsg(p, a1, a2, a3)
	case common.SYS_GETSOCKOPT:
		ret = sys_getsockopt(p, a1, a2, a3, a4, a5)
	case common.SYS_SETSOCKOPT:
		ret = sys_setsockopt(p, a1, a2, a3, a4, a5)
	case common.SYS_FORK:
		ret = sys_fork(p, tf, a1, a2)
	case common.SYS_EXECV:
		ret = sys_execv(p, tf, a1, a2)
	case common.SYS_EXIT:
		status := a1 & 0xff
		status |= common.EXITED
		s.Sys_exit(p, tid, status)
	case common.SYS_WAIT4:
		ret = sys_wait4(p, tid, a1, a2, a3, a4, a5)
	case common.SYS_KILL:
		ret = sys_kill(p, a1, a2)
	case common.SYS_FCNTL:
		ret = sys_fcntl(p, a1, a2, a3)
	case common.SYS_TRUNC:
		ret = sys_truncate(p, a1, uint(a2))
	case common.SYS_FTRUNC:
		ret = sys_ftruncate(p, a1, uint(a2))
	case common.SYS_GETCWD:
		ret = sys_getcwd(p, a1, a2)
	case common.SYS_CHDIR:
		ret = sys_chdir(p, a1)
	case common.SYS_RENAME:
		ret = sys_rename(p, a1, a2)
	case common.SYS_MKDIR:
		ret = sys_mkdir(p, a1, a2)
	case common.SYS_LINK:
		ret = sys_link(p, a1, a2)
	case common.SYS_UNLINK:
		ret = sys_unlink(p, a1, a2)
	case common.SYS_GETTOD:
		ret = sys_gettimeofday(p, a1)
	case common.SYS_GETRLMT:
		ret = sys_getrlimit(p, a1, a2)
	case common.SYS_GETRUSG:
		ret = sys_getrusage(p, a1, a2)
	case common.SYS_MKNOD:
		ret = sys_mknod(p, a1, a2, a3)
	case common.SYS_SETRLMT:
		ret = sys_setrlimit(p, a1, a2)
	case common.SYS_SYNC:
		ret = sys_sync(p)
	case common.SYS_REBOOT:
		ret = sys_reboot(p)
	case common.SYS_NANOSLEEP:
		ret = sys_nanosleep(p, a1, a2)
	case common.SYS_PIPE2:
		ret = sys_pipe2(p, a1, a2)
	case common.SYS_PROF:
		ret = sys_prof(p, a1, a2, a3, a4)
	case common.SYS_INFO:
		ret = sys_info(p, a1)
	case common.SYS_THREXIT:
		sys_threxit(p, tid, a1)
	case common.SYS_PREAD:
		ret = sys_pread(p, a1, a2, a3, a4)
	case common.SYS_PWRITE:
		ret = sys_pwrite(p, a1, a2, a3, a4)
	case common.SYS_FUTEX:
		ret = sys_futex(p, a1, a2, a3, a4, a5)
	case common.SYS_GETTID:
		ret = sys_gettid(p, tid)
	default:
		fmt.Printf("unexpected syscall %v\n", sysno)
		s.Sys_exit(p, tid, common.SIGNALED|common.Mkexitsig(31))
	}
	return ret
}

// Implements Console_i
type console_t struct {
}

var console = &console_t{}

func (c *console_t) Cons_read(ub vm.Userio_i, offset int) (int, defs.Err_t) {
	sz := ub.Remain()
	kdata, err := kbd_get(sz)
	if err != 0 {
		return 0, err
	}
	ret, err := ub.Uiowrite(kdata)
	if err != 0 || ret != len(kdata) {
		fmt.Printf("dropped keys!\n")
	}
	return ret, err
}

func (c *console_t) Cons_write(src vm.Userio_i, off int) (int, defs.Err_t) {
	// merge into one buffer to avoid taking the console lock many times.
	// what a sweet optimization.
	utext := int8(0x17)
	big := make([]uint8, src.Totalsz())
	read, err := src.Uioread(big)
	if err != 0 {
		return 0, err
	}
	if read != src.Totalsz() {
		panic("short read")
	}
	runtime.Pmsga(&big[0], len(big), utext)
	return len(big), 0
}

func _fd_read(proc *common.Proc_t, fdn int) (*vm.Fd_t, defs.Err_t) {
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return nil, -defs.EBADF
	}
	if fd.Perms&vm.FD_READ == 0 {
		return nil, -defs.EPERM
	}
	return fd, 0
}

func _fd_write(proc *common.Proc_t, fdn int) (*vm.Fd_t, defs.Err_t) {
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return nil, -defs.EBADF
	}
	if fd.Perms&vm.FD_WRITE == 0 {
		return nil, -defs.EPERM
	}
	return fd, 0
}

func sys_read(proc *common.Proc_t, fdn int, bufp int, sz int) int {
	if sz == 0 {
		return 0
	}
	fd, err := _fd_read(proc, fdn)
	if err != 0 {
		return int(err)
	}
	userbuf := proc.Aspace.Mkuserbuf(bufp, sz)

	ret, err := fd.Fops.Read(userbuf)
	if err != 0 {
		return int(err)
	}
	//common.Ubpool.Put(userbuf)
	return ret
}

func sys_write(proc *common.Proc_t, fdn int, bufp int, sz int) int {
	if sz == 0 {
		return 0
	}
	fd, err := _fd_write(proc, fdn)
	if err != 0 {
		return int(err)
	}
	userbuf := proc.Aspace.Mkuserbuf(bufp, sz)

	ret, err := fd.Fops.Write(userbuf)
	if err != 0 {
		return int(err)
	}
	//common.Ubpool.Put(userbuf)
	return ret
}

func sys_open(proc *common.Proc_t, pathn int, _flags int, mode int) int {
	path, err := proc.Aspace.Userstr(pathn, fs.NAME_MAX)
	if err != 0 {
		return int(err)
	}
	flags := common.Fdopt_t(_flags)
	temp := flags & (common.O_RDONLY | common.O_WRONLY | common.O_RDWR)
	if temp != common.O_RDONLY && temp != common.O_WRONLY && temp != common.O_RDWR {
		return int(-defs.EINVAL)
	}
	if temp == common.O_RDONLY && flags&common.O_TRUNC != 0 {
		return int(-defs.EINVAL)
	}
	fdperms := 0
	switch temp {
	case common.O_RDONLY:
		fdperms = vm.FD_READ
	case common.O_WRONLY:
		fdperms = vm.FD_WRITE
	case common.O_RDWR:
		fdperms = vm.FD_READ | vm.FD_WRITE
	default:
		fdperms = vm.FD_READ
	}
	err = badpath(path)
	if err != 0 {
		return int(err)
	}
	file, err := thefs.Fs_open(path, flags, mode, proc.Cwd, 0, 0)
	if err != 0 {
		return int(err)
	}
	if flags&common.O_CLOEXEC != 0 {
		fdperms |= vm.FD_CLOEXEC
	}
	fdn, ok := proc.Fd_insert(file, fdperms)
	if !ok {
		lhits++
		vm.Close_panic(file)
		return int(-defs.EMFILE)
	}
	return fdn
}

func sys_pause(proc *common.Proc_t) int {
	// no signals yet!
	var c chan bool
	select {
	case <-c:
	case <-tinfo.Current().Killnaps.Killch:
	}
	return -1
}

func (s *syscall_t) Sys_close(proc *common.Proc_t, fdn int) int {
	fd, ok := proc.Fd_del(fdn)
	if !ok {
		return int(-defs.EBADF)
	}
	ret := fd.Fops.Close()
	return int(ret)
}

func sys_mmap(proc *common.Proc_t, addrn, lenn, protflags, fdn, offset int) int {
	if lenn == 0 {
		return int(-defs.EINVAL)
	}
	prot := uint(protflags) >> 32
	flags := uint(uint32(protflags))

	mask := common.MAP_SHARED | common.MAP_PRIVATE
	if flags&mask == 0 || flags&mask == mask {
		return int(-defs.EINVAL)
	}
	shared := flags&common.MAP_SHARED != 0
	anon := flags&common.MAP_ANON != 0
	fdmap := !anon
	if (fdmap && fdn < 0) || (fdmap && offset < 0) || (anon && fdn >= 0) {
		return int(-defs.EINVAL)
	}
	if flags&common.MAP_FIXED != 0 {
		return int(-defs.EINVAL)
	}
	// OpenBSD allows mappings of only PROT_WRITE and read accesses that
	// fault-in the page cause a segfault while writes do not. Reads
	// following a write do not cause segfault (of course). POSIX
	// apparently requires an implementation to support only common.PROT_WRITE,
	// but it seems better to disallow permission schemes that the CPU
	// cannot enforce.
	if prot&common.PROT_READ == 0 {
		return int(-defs.EINVAL)
	}
	if prot == common.PROT_NONE {
		panic("no imp")
		return proc.Mmapi
	}

	var fd *vm.Fd_t
	if fdmap {
		var ok bool
		fd, ok = proc.Fd_get(fdn)
		if !ok {
			return int(-defs.EBADF)
		}
		if fd.Perms&vm.FD_READ == 0 ||
			(shared && prot&common.PROT_WRITE != 0 &&
				fd.Perms&vm.FD_WRITE == 0) {
			return int(-defs.EACCES)
		}
	}

	proc.Aspace.Lock_pmap()

	perms := vm.PTE_U
	if prot&common.PROT_WRITE != 0 {
		perms |= vm.PTE_W
	}
	lenn = util.Roundup(lenn, mem.PGSIZE)
	// limit checks
	if lenn/int(mem.PGSIZE)+proc.Aspace.Vmregion.Pglen() > proc.Ulim.Pages {
		proc.Aspace.Unlock_pmap()
		lhits++
		return int(-defs.ENOMEM)
	}
	if proc.Aspace.Vmregion.Novma >= proc.Ulim.Novma {
		proc.Aspace.Unlock_pmap()
		lhits++
		return int(-defs.ENOMEM)
	}

	addr := proc.Aspace.Unusedva_inner(proc.Mmapi, lenn)
	proc.Mmapi = addr + lenn
	switch {
	case anon && shared:
		proc.Aspace.Vmadd_shareanon(addr, lenn, perms)
	case anon && !shared:
		proc.Aspace.Vmadd_anon(addr, lenn, perms)
	case fdmap:
		fops := fd.Fops
		// vmadd_*file will increase the open count on the file
		if shared {
			proc.Aspace.Vmadd_sharefile(addr, lenn, perms, fops, offset,
				thefs)
		} else {
			proc.Aspace.Vmadd_file(addr, lenn, perms, fops, offset)
		}
	}
	tshoot := false
	// eagerly map anonymous pages, lazily-map file pages. our vm system
	// supports lazily-mapped private anonymous pages though.
	var ub int
	failed := false
	if anon {
		for i := 0; i < lenn; i += int(mem.PGSIZE) {
			_, p_pg, ok := physmem.Refpg_new()
			if !ok {
				failed = true
				break
			}
			ns, ok := proc.Aspace.Page_insert(addr+i, p_pg, perms, true)
			if !ok {
				physmem.Refdown(p_pg)
				failed = true
				break
			}
			ub = i
			tshoot = tshoot || ns
		}
	}
	ret := addr
	if failed {
		for i := 0; i < ub; i += mem.PGSIZE {
			proc.Aspace.Page_remove(addr + i)
		}
		// removing this region cannot create any more vm objects than
		// what this call to sys_mmap started with.
		if proc.Aspace.Vmregion.Remove(addr, lenn, proc.Ulim.Novma) != 0 {
			panic("wut")
		}
		ret = int(-defs.ENOMEM)
	}
	// sys_mmap won't replace pages since it always finds unused VA space,
	// so the following TLB shootdown is never used.
	if tshoot {
		proc.Aspace.Tlbshoot(0, 1)
	}
	proc.Aspace.Unlock_pmap()
	return ret
}

func sys_munmap(proc *common.Proc_t, addrn, len int) int {
	if addrn&int(vm.PGOFFSET) != 0 || addrn < mem.USERMIN {
		return int(-defs.EINVAL)
	}
	proc.Aspace.Lock_pmap()
	defer proc.Aspace.Unlock_pmap()

	vmi1, ok1 := proc.Aspace.Vmregion.Lookup(uintptr(addrn))
	vmi2, ok2 := proc.Aspace.Vmregion.Lookup(uintptr(addrn+len) - 1)
	if !ok1 || !ok2 || vmi1.Pgn != vmi2.Pgn {
		return int(-defs.EINVAL)
	}

	err := proc.Aspace.Vmregion.Remove(addrn, len, proc.Ulim.Novma)
	if err != 0 {
		lhits++
		return int(err)
	}
	// addrn must be page-aligned
	len = util.Roundup(len, mem.PGSIZE)
	for i := 0; i < len; i += mem.PGSIZE {
		p := addrn + i
		if p < mem.USERMIN {
			panic("how")
		}
		proc.Aspace.Page_remove(p)
	}
	pgs := len >> vm.PGSHIFT
	proc.Aspace.Tlbshoot(uintptr(addrn), pgs)
	return 0
}

func sys_readv(proc *common.Proc_t, fdn, _iovn, iovcnt int) int {
	fd, err := _fd_read(proc, fdn)
	if err != 0 {
		return int(err)
	}
	iovn := uint(_iovn)
	iov := &vm.Useriovec_t{}
	if err := iov.Iov_init(&proc.Aspace, iovn, iovcnt); err != 0 {
		return int(err)
	}
	ret, err := fd.Fops.Read(iov)
	if err != 0 {
		return int(err)
	}
	return ret
}

func sys_writev(proc *common.Proc_t, fdn, _iovn, iovcnt int) int {
	fd, err := _fd_write(proc, fdn)
	if err != 0 {
		return int(err)
	}
	iovn := uint(_iovn)
	iov := &vm.Useriovec_t{}
	if err := iov.Iov_init(&proc.Aspace, iovn, iovcnt); err != 0 {
		return int(err)
	}
	ret, err := fd.Fops.Write(iov)
	if err != 0 {
		return int(err)
	}
	return ret
}

func sys_sigaction(proc *common.Proc_t, sig, actn, oactn int) int {
	panic("no imp")
}

func sys_access(proc *common.Proc_t, pathn, mode int) int {
	path, err := proc.Aspace.Userstr(pathn, fs.NAME_MAX)
	if err != 0 {
		return int(err)
	}
	if mode == 0 {
		return int(-defs.EINVAL)
	}

	fsf, err := thefs.Fs_open_inner(path, common.O_RDONLY, 0, proc.Cwd, 0, 0)
	if err != 0 {
		return int(err)
	}

	// XXX no permissions yet
	//R_OK := 1 << 0
	//W_OK := 1 << 1
	//X_OK := 1 << 2
	ret := 0

	if thefs.Fs_close(fsf.Inum) != 0 {
		panic("must succeed")
	}
	return ret
}

func sys_dup2(proc *common.Proc_t, oldn, newn int) int {
	if oldn == newn {
		return newn
	}
	ofd, needclose, err := proc.Fd_dup(oldn, newn)
	if err != 0 {
		return int(err)
	}
	if needclose {
		vm.Close_panic(ofd)
	}
	return newn
}

func sys_stat(proc *common.Proc_t, pathn, statn int) int {
	path, err := proc.Aspace.Userstr(pathn, fs.NAME_MAX)
	if err != 0 {
		return int(err)
	}
	buf := &stat.Stat_t{}
	err = thefs.Fs_stat(path, buf, proc.Cwd)
	if err != 0 {
		return int(err)
	}
	return int(proc.Aspace.K2user(buf.Bytes(), statn))
}

func sys_fstat(proc *common.Proc_t, fdn int, statn int) int {
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return int(-defs.EBADF)
	}
	buf := &stat.Stat_t{}
	err := fd.Fops.Fstat(buf)
	if err != 0 {
		return int(err)
	}

	return int(proc.Aspace.K2user(buf.Bytes(), statn))
}

// converts internal states to poll states
// pokes poll status bits into user memory. since we only use one priority
// internally, mask away any POLL bits the user didn't not request.
func _ready2rev(orig int, r vm.Ready_t) int {
	inmask := common.POLLIN | common.POLLPRI
	outmask := common.POLLOUT | common.POLLWRBAND
	pbits := 0
	if r&vm.R_READ != 0 {
		pbits |= inmask
	}
	if r&vm.R_WRITE != 0 {
		pbits |= outmask
	}
	if r&vm.R_HUP != 0 {
		pbits |= common.POLLHUP
	}
	if r&vm.R_ERROR != 0 {
		pbits |= common.POLLERR
	}
	wantevents := ((orig >> 32) & 0xffff) | common.POLLNVAL | common.POLLERR | common.POLLHUP
	revents := wantevents & pbits
	return orig | (revents << 48)
}

func _checkfds(proc *common.Proc_t, tid defs.Tid_t, pm *vm.Pollmsg_t, wait bool, buf []uint8,
	nfds int) (int, bool, defs.Err_t) {
	inmask := common.POLLIN | common.POLLPRI
	outmask := common.POLLOUT | common.POLLWRBAND
	readyfds := 0
	writeback := false
	proc.Fdl.Lock()
	for i := 0; i < nfds; i++ {
		off := i * 8
		uw := readn(buf, 8, off)
		fdn := int(uint32(uw))
		// fds < 0 are to be ignored
		if fdn < 0 {
			continue
		}
		fd, ok := proc.Fd_get_inner(fdn)
		if !ok {
			uw |= common.POLLNVAL
			writen(buf, 8, off, uw)
			writeback = true
			continue
		}
		var pev vm.Ready_t
		events := int((uint(uw) >> 32) & 0xffff)
		// one priority
		if events&inmask != 0 {
			pev |= vm.R_READ
		}
		if events&outmask != 0 {
			pev |= vm.R_WRITE
		}
		if events&common.POLLHUP != 0 {
			pev |= vm.R_HUP
		}
		// poll unconditionally reports ERR, HUP, and NVAL
		pev |= vm.R_ERROR | vm.R_HUP
		pm.Pm_set(tid, pev, wait)
		devstatus, err := fd.Fops.Pollone(*pm)
		if err != 0 {
			proc.Fdl.Unlock()
			return 0, false, err
		}
		if devstatus != 0 {
			// found at least one ready fd; don't bother having the
			// other fds send notifications. update user revents
			wait = false
			nuw := _ready2rev(uw, devstatus)
			writen(buf, 8, off, nuw)
			readyfds++
			writeback = true
		}
	}
	proc.Fdl.Unlock()
	return readyfds, writeback, 0
}

func sys_poll(proc *common.Proc_t, tid defs.Tid_t, fdsn, nfds, timeout int) int {
	if nfds < 0 || timeout < -1 {
		return int(-defs.EINVAL)
	}

	// copy pollfds from userspace to avoid reading/writing overhead
	// (locking pmap and looking up uva mapping).
	pollfdsz := 8
	sz := uint(pollfdsz * nfds)
	// chosen arbitrarily...
	maxsz := uint(4096)
	if sz > maxsz {
		// fall back to holding lock over user pmap if they want to
		// poll so many fds.
		fmt.Printf("poll limit hit\n")
		return int(-defs.EINVAL)
	}
	buf := make([]uint8, sz)
	if err := proc.Aspace.User2k(buf, fdsn); err != 0 {
		return int(err)
	}

	// first we tell the underlying device to notify us if their fd is
	// ready. if a device is immediately ready, we don't bother to register
	// notifiers with the rest of the devices -- we just ask their status
	// too.
	gimme := bounds.Bounds(bounds.B_SYS_POLL)
	pm := vm.Pollmsg_t{}
	for {
		// its ok to block for memory here since no locks are held
		if !res.Resadd(gimme) {
			return int(-defs.ENOHEAP)
		}
		wait := timeout != 0
		rfds, writeback, err := _checkfds(proc, tid, &pm, wait, buf,
			nfds)
		if err != 0 {
			return int(err)
		}
		if writeback {
			if err := proc.Aspace.K2user(buf, fdsn); err != 0 {
				return int(err)
			}
		}

		// if we found a ready fd, we are done
		if rfds != 0 || !wait {
			return rfds
		}

		// otherwise, wait for a notification
		timedout, err := pm.Pm_wait(timeout)
		if err != 0 {
			return int(err)
		}
		if timedout {
			return 0
		}
	}
}

func sys_lseek(proc *common.Proc_t, fdn, off, whence int) int {
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return int(-defs.EBADF)
	}

	ret, err := fd.Fops.Lseek(off, whence)
	if err != 0 {
		return int(err)
	}
	return ret
}

func sys_pipe2(proc *common.Proc_t, pipen, _flags int) int {
	rfp := vm.FD_READ
	wfp := vm.FD_WRITE

	flags := common.Fdopt_t(_flags)
	var opts common.Fdopt_t
	if flags&common.O_NONBLOCK != 0 {
		opts |= common.O_NONBLOCK
	}

	if flags&common.O_CLOEXEC != 0 {
		rfp |= vm.FD_CLOEXEC
		wfp |= vm.FD_CLOEXEC
	}

	// if there is an error, pipe_t.op_reopen() will release the pipe
	// reservation.
	if !limits.Syslimit.Pipes.Take() {
		lhits++
		return int(-defs.ENOMEM)
	}

	p := &pipe_t{lraise: true}
	p.pipe_start()
	rops := &pipefops_t{pipe: p, writer: false, options: opts}
	wops := &pipefops_t{pipe: p, writer: true, options: opts}
	rpipe := &vm.Fd_t{Fops: rops}
	wpipe := &vm.Fd_t{Fops: wops}
	rfd, wfd, ok := proc.Fd_insert2(rpipe, rfp, wpipe, wfp)
	if !ok {
		vm.Close_panic(rpipe)
		vm.Close_panic(wpipe)
		return int(-defs.EMFILE)
	}

	err := proc.Aspace.Userwriten(pipen, 4, rfd)
	if err != 0 {
		goto bail
	}
	err = proc.Aspace.Userwriten(pipen+4, 4, wfd)
	if err != 0 {
		goto bail
	}
	return 0

bail:
	err1 := sys.Sys_close(proc, rfd)
	err2 := sys.Sys_close(proc, wfd)
	if err1 != 0 || err2 != 0 {
		panic("must succeed")
	}
	return int(err)
}

type pipe_t struct {
	sync.Mutex
	cbuf    circbuf_t
	rcond   *sync.Cond
	wcond   *sync.Cond
	readers int
	writers int
	closed  bool
	pollers vm.Pollers_t
	passfds passfd_t
	// if true, this pipe was allocated against the pipe limit; raise it on
	// termination.
	lraise bool
}

func (o *pipe_t) pipe_start() {
	pipesz := mem.PGSIZE
	o.cbuf.cb_init(pipesz)
	o.readers, o.writers = 1, 1
	o.rcond = sync.NewCond(o)
	o.wcond = sync.NewCond(o)
}

func (o *pipe_t) op_write(src vm.Userio_i, noblock bool) (int, defs.Err_t) {
	const pipe_buf = 4096
	need := src.Remain()
	if need > pipe_buf {
		if noblock {
			need = 1
		} else {
			need = pipe_buf
		}
	}
	o.Lock()
	for {
		if o.closed {
			o.Unlock()
			return 0, -defs.EBADF
		}
		if o.readers == 0 {
			o.Unlock()
			return 0, -defs.EPIPE
		}
		if o.cbuf.left() >= need {
			break
		}
		if noblock {
			o.Unlock()
			return 0, -defs.EWOULDBLOCK
		}
		if err := common.KillableWait(o.wcond); err != 0 {
			o.Unlock()
			return 0, err
		}
	}
	ret, err := o.cbuf.copyin(src)
	if err != 0 {
		o.Unlock()
		return 0, err
	}
	o.rcond.Signal()
	o.pollers.Wakeready(vm.R_READ)
	o.Unlock()

	return ret, 0
}

func (o *pipe_t) op_read(dst vm.Userio_i, noblock bool) (int, defs.Err_t) {
	o.Lock()
	for {
		if o.closed {
			o.Unlock()
			return 0, -defs.EBADF
		}
		if o.writers == 0 || !o.cbuf.empty() {
			break
		}
		if noblock {
			o.Unlock()
			return 0, -defs.EWOULDBLOCK
		}
		if err := common.KillableWait(o.rcond); err != 0 {
			o.Unlock()
			return 0, err
		}
	}
	ret, err := o.cbuf.copyout(dst)
	if err != 0 {
		o.Unlock()
		return 0, err
	}
	o.wcond.Signal()
	o.pollers.Wakeready(vm.R_WRITE)
	o.Unlock()

	return ret, 0
}

func (o *pipe_t) op_poll(pm vm.Pollmsg_t) (vm.Ready_t, defs.Err_t) {
	o.Lock()

	if o.closed {
		o.Unlock()
		return 0, 0
	}

	var r vm.Ready_t
	readable := false
	if !o.cbuf.empty() || o.writers == 0 {
		readable = true
	}
	writeable := false
	if !o.cbuf.full() || o.readers == 0 {
		writeable = true
	}
	if pm.Events&vm.R_READ != 0 && readable {
		r |= vm.R_READ
	}
	if pm.Events&vm.R_HUP != 0 && o.writers == 0 {
		r |= vm.R_HUP
	} else if pm.Events&vm.R_WRITE != 0 && writeable {
		r |= vm.R_WRITE
	}
	if r != 0 || !pm.Dowait {
		o.Unlock()
		return r, 0
	}
	err := o.pollers.Addpoller(&pm)
	o.Unlock()
	return 0, err
}

func (o *pipe_t) op_reopen(rd, wd int) defs.Err_t {
	o.Lock()
	if o.closed {
		o.Unlock()
		return -defs.EBADF
	}
	o.readers += rd
	o.writers += wd
	if o.writers == 0 {
		o.rcond.Broadcast()
	}
	if o.readers == 0 {
		o.wcond.Broadcast()
	}
	if o.readers == 0 && o.writers == 0 {
		o.closed = true
		o.cbuf.cb_release()
		o.passfds.closeall()
		if o.lraise {
			limits.Syslimit.Pipes.Give()
		}
	}
	o.Unlock()
	return 0
}

func (o *pipe_t) op_fdadd(nfd *vm.Fd_t) defs.Err_t {
	o.Lock()
	defer o.Unlock()

	for !o.passfds.add(nfd) {
		if err := common.KillableWait(o.wcond); err != 0 {
			return err
		}
	}
	return 0
}

func (o *pipe_t) op_fdtake() (*vm.Fd_t, bool) {
	o.Lock()
	defer o.Unlock()
	ret, ok := o.passfds.take()
	if !ok {
		return nil, false
	}
	o.wcond.Broadcast()
	return ret, true
}

type pipefops_t struct {
	pipe    *pipe_t
	options common.Fdopt_t
	writer  bool
}

func (of *pipefops_t) Close() defs.Err_t {
	var ret defs.Err_t
	if of.writer {
		ret = of.pipe.op_reopen(0, -1)
	} else {
		ret = of.pipe.op_reopen(-1, 0)
	}
	return ret
}

func (of *pipefops_t) Fstat(st *stat.Stat_t) defs.Err_t {
	// linux and openbsd give same mode for all pipes
	st.Wdev(0)
	pipemode := uint(3 << 16)
	st.Wmode(pipemode)
	return 0
}

func (of *pipefops_t) Lseek(int, int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

func (of *pipefops_t) Mmapi(int, int, bool) ([]mem.Mmapinfo_t, defs.Err_t) {
	return nil, -defs.EINVAL
}

func (of *pipefops_t) Pathi() defs.Inum_t {
	panic("pipe cwd")
}

func (of *pipefops_t) Read(dst vm.Userio_i) (int, defs.Err_t) {
	noblk := of.options&common.O_NONBLOCK != 0
	return of.pipe.op_read(dst, noblk)
}

func (of *pipefops_t) Reopen() defs.Err_t {
	var ret defs.Err_t
	if of.writer {
		ret = of.pipe.op_reopen(0, 1)
	} else {
		ret = of.pipe.op_reopen(1, 0)
	}
	return ret
}

func (of *pipefops_t) Write(src vm.Userio_i) (int, defs.Err_t) {
	noblk := of.options&common.O_NONBLOCK != 0
	c := 0
	for c != src.Totalsz() {
		if !res.Resadd(bounds.Bounds(bounds.B_PIPEFOPS_T_WRITE)) {
			return c, -defs.ENOHEAP
		}
		ret, err := of.pipe.op_write(src, noblk)
		if noblk || err != 0 {
			return ret, err
		}
		c += ret
	}
	return c, 0
}

func (of *pipefops_t) Truncate(uint) defs.Err_t {
	return -defs.EINVAL
}

func (of *pipefops_t) Pread(vm.Userio_i, int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

func (of *pipefops_t) Pwrite(vm.Userio_i, int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

func (of *pipefops_t) Accept(vm.Userio_i) (vm.Fdops_i, int, defs.Err_t) {
	return nil, 0, -defs.ENOTSOCK
}

func (of *pipefops_t) Bind([]uint8) defs.Err_t {
	return -defs.ENOTSOCK
}

func (of *pipefops_t) Connect([]uint8) defs.Err_t {
	return -defs.ENOTSOCK
}

func (of *pipefops_t) Listen(int) (vm.Fdops_i, defs.Err_t) {
	return nil, -defs.ENOTSOCK
}

func (of *pipefops_t) Sendmsg(vm.Userio_i, []uint8, []uint8,
	int) (int, defs.Err_t) {
	return 0, -defs.ENOTSOCK
}

func (of *pipefops_t) Recvmsg(vm.Userio_i, vm.Userio_i,
	vm.Userio_i, int) (int, int, int, defs.Msgfl_t, defs.Err_t) {
	return 0, 0, 0, 0, -defs.ENOTSOCK
}

func (of *pipefops_t) Pollone(pm vm.Pollmsg_t) (vm.Ready_t, defs.Err_t) {
	if of.writer {
		pm.Events &^= vm.R_READ
	} else {
		pm.Events &^= vm.R_WRITE
	}
	return of.pipe.op_poll(pm)
}

func (of *pipefops_t) Fcntl(cmd, opt int) int {
	switch cmd {
	case common.F_GETFL:
		return int(of.options)
	case common.F_SETFL:
		of.options = common.Fdopt_t(opt)
		return 0
	default:
		panic("weird cmd")
	}
}

func (of *pipefops_t) Getsockopt(int, vm.Userio_i, int) (int, defs.Err_t) {
	return 0, -defs.ENOTSOCK
}

func (of *pipefops_t) Setsockopt(int, int, vm.Userio_i, int) defs.Err_t {
	return -defs.ENOTSOCK
}

func (of *pipefops_t) Shutdown(read, write bool) defs.Err_t {
	return -defs.ENOTCONN
}

func sys_rename(proc *common.Proc_t, oldn int, newn int) int {
	old, err1 := proc.Aspace.Userstr(oldn, fs.NAME_MAX)
	new, err2 := proc.Aspace.Userstr(newn, fs.NAME_MAX)
	if err1 != 0 {
		return int(err1)
	}
	if err2 != 0 {
		return int(err2)
	}
	err1 = badpath(old)
	err2 = badpath(new)
	if err1 != 0 {
		return int(err1)
	}
	if err2 != 0 {
		return int(err2)
	}
	err := thefs.Fs_rename(old, new, proc.Cwd)
	return int(err)
}

func sys_mkdir(proc *common.Proc_t, pathn int, mode int) int {
	path, err := proc.Aspace.Userstr(pathn, fs.NAME_MAX)
	if err != 0 {
		return int(err)
	}
	err = badpath(path)
	if err != 0 {
		return int(err)
	}
	err = thefs.Fs_mkdir(path, mode, proc.Cwd)
	return int(err)
}

func sys_link(proc *common.Proc_t, oldn int, newn int) int {
	old, err1 := proc.Aspace.Userstr(oldn, fs.NAME_MAX)
	new, err2 := proc.Aspace.Userstr(newn, fs.NAME_MAX)
	if err1 != 0 {
		return int(err1)
	}
	if err2 != 0 {
		return int(err2)
	}
	err1 = badpath(old)
	err2 = badpath(new)
	if err1 != 0 {
		return int(err1)
	}
	if err2 != 0 {
		return int(err2)
	}
	err := thefs.Fs_link(old, new, proc.Cwd)
	return int(err)
}

func sys_unlink(proc *common.Proc_t, pathn, isdiri int) int {
	path, err := proc.Aspace.Userstr(pathn, fs.NAME_MAX)
	if err != 0 {
		return int(err)
	}
	err = badpath(path)
	if err != 0 {
		return int(err)
	}
	wantdir := isdiri != 0
	err = thefs.Fs_unlink(path, proc.Cwd, wantdir)
	return int(err)
}

func sys_gettimeofday(proc *common.Proc_t, timevaln int) int {
	tvalsz := 16
	now := time.Now()
	buf := make([]uint8, tvalsz)
	us := int(now.UnixNano() / 1000)
	writen(buf, 8, 0, us/1e6)
	writen(buf, 8, 8, us%1e6)
	if err := proc.Aspace.K2user(buf, timevaln); err != 0 {
		return int(err)
	}
	return 0
}

var _rlimits = map[int]uint{common.RLIMIT_NOFILE: common.RLIM_INFINITY}

func sys_getrlimit(proc *common.Proc_t, resn, rlpn int) int {
	var cur uint
	switch resn {
	case common.RLIMIT_NOFILE:
		cur = proc.Ulim.Nofile
	default:
		return int(-defs.EINVAL)
	}
	max := _rlimits[resn]
	err1 := proc.Aspace.Userwriten(rlpn, 8, int(cur))
	err2 := proc.Aspace.Userwriten(rlpn+8, 8, int(max))
	if err1 != 0 {
		return int(err1)
	}
	if err2 != 0 {
		return int(err2)
	}
	return 0
}

func sys_setrlimit(proc *common.Proc_t, resn, rlpn int) int {
	// XXX root can raise max
	_ncur, err := proc.Aspace.Userreadn(rlpn, 8)
	if err != 0 {
		return int(err)
	}
	ncur := uint(_ncur)
	if ncur > _rlimits[resn] {
		return int(-defs.EINVAL)
	}
	switch resn {
	case common.RLIMIT_NOFILE:
		proc.Ulim.Nofile = ncur
	default:
		return int(-defs.EINVAL)
	}
	return 0
}

func sys_getrusage(proc *common.Proc_t, who, rusagep int) int {
	var ru []uint8
	if who == common.RUSAGE_SELF {
		// user time is gathered at thread termination... report user
		// time as best as we can
		tmp := proc.Atime

		proc.Threadi.Lock()
		for tid := range proc.Threadi.Notes {
			if tid == 0 {
			}
			val := 42
			// tid may not exist if the query for the time races
			// with a thread exiting.
			if val > 0 {
				tmp.Userns += int64(val)
			}
		}
		proc.Threadi.Unlock()

		ru = tmp.To_rusage()
	} else if who == common.RUSAGE_CHILDREN {
		ru = proc.Catime.Fetch()
	} else {
		return int(-defs.EINVAL)
	}
	if err := proc.Aspace.K2user(ru, rusagep); err != 0 {
		return int(err)
	}
	return int(-defs.ENOSYS)
}

func mkdev(_maj, _min int) uint {
	maj := uint(_maj)
	min := uint(_min)
	if min > 0xff {
		panic("bad minor")
	}
	m := maj<<8 | min
	return uint(m << 32)
}

func unmkdev(d uint) (int, int) {
	return int(d >> 40), int(uint8(d >> 32))
}

func sys_mknod(proc *common.Proc_t, pathn, moden, devn int) int {
	path, err := proc.Aspace.Userstr(pathn, fs.NAME_MAX)
	if err != 0 {
		return int(err)
	}

	err = badpath(path)
	if err != 0 {
		return int(err)
	}
	maj, min := unmkdev(uint(devn))
	fsf, err := thefs.Fs_open_inner(path, common.O_CREAT, 0, proc.Cwd, maj, min)
	if err != 0 {
		return int(err)
	}
	if thefs.Fs_close(fsf.Inum) != 0 {
		panic("must succeed")
	}
	return 0
}

func sys_sync(proc *common.Proc_t) int {
	return int(thefs.Fs_sync())
}

func sys_reboot(proc *common.Proc_t) int {
	// mov'ing to cr3 does not flush global pages. if, before loading the
	// zero page into cr3 below, there are just enough TLB entries to
	// dispatch a fault, but not enough to complete the fault handler, the
	// fault handler will recursively fault forever since it uses an IST
	// stack. therefore, flush the global pages too.
	pge := uintptr(1 << 7)
	runtime.Lcr4(runtime.Rcr4() &^ pge)
	// who needs ACPI?
	runtime.Lcr3(uintptr(mem.P_zeropg))
	// poof
	fmt.Printf("what?\n")
	return 0
}

func sys_nanosleep(proc *common.Proc_t, sleeptsn, remaintsn int) int {
	tot, _, err := proc.Aspace.Usertimespec(sleeptsn)
	if err != 0 {
		return int(err)
	}
	tochan := time.After(tot)
	kn := &tinfo.Current().Killnaps
	select {
	case <-tochan:
		return 0
	case <-kn.Killch:
		if kn.Kerr == 0 {
			panic("no")
		}
		return int(kn.Kerr)
	}
}

func sys_getpid(proc *common.Proc_t, tid defs.Tid_t) int {
	return proc.Pid
}

func sys_getppid(proc *common.Proc_t, tid defs.Tid_t) int {
	return proc.Pwait.Pid
}

func sys_socket(proc *common.Proc_t, domain, typ, proto int) int {
	var opts common.Fdopt_t
	if typ&common.SOCK_NONBLOCK != 0 {
		opts |= common.O_NONBLOCK
	}
	var clop int
	if typ&common.SOCK_CLOEXEC != 0 {
		clop = vm.FD_CLOEXEC
	}

	var sfops vm.Fdops_i
	switch {
	case domain == common.AF_UNIX && typ&common.SOCK_DGRAM != 0:
		if opts != 0 {
			panic("no imp")
		}
		sfops = &sudfops_t{open: 1}
	case domain == common.AF_UNIX && typ&common.SOCK_STREAM != 0:
		sfops = &susfops_t{options: opts}
	case domain == common.AF_INET && typ&common.SOCK_STREAM != 0:
		tfops := &tcpfops_t{tcb: &tcptcb_t{}, options: opts}
		tfops.tcb.openc = 1
		sfops = tfops
	default:
		return int(-defs.EINVAL)
	}
	if !limits.Syslimit.Socks.Take() {
		lhits++
		return int(-defs.ENOMEM)
	}
	file := &vm.Fd_t{}
	file.Fops = sfops
	fdn, ok := proc.Fd_insert(file, vm.FD_READ|vm.FD_WRITE|clop)
	if !ok {
		vm.Close_panic(file)
		limits.Syslimit.Socks.Give()
		return int(-defs.EMFILE)
	}
	return fdn
}

func sys_connect(proc *common.Proc_t, fdn, sockaddrn, socklen int) int {
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return int(-defs.EBADF)
	}

	// copy sockaddr to kernel space to avoid races
	sabuf, err := copysockaddr(proc, sockaddrn, socklen)
	if err != 0 {
		return int(err)
	}
	err = fd.Fops.Connect(sabuf)
	return int(err)
}

func sys_accept(proc *common.Proc_t, fdn, sockaddrn, socklenn int) int {
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return int(-defs.EBADF)
	}
	var sl int
	if socklenn != 0 {
		l, err := proc.Aspace.Userreadn(socklenn, 8)
		if err != 0 {
			return int(err)
		}
		if l < 0 {
			return int(-defs.EFAULT)
		}
		sl = l
	}
	fromsa := proc.Aspace.Mkuserbuf(sockaddrn, sl)
	newfops, fromlen, err := fd.Fops.Accept(fromsa)
	//common.Ubpool.Put(fromsa)
	if err != 0 {
		return int(err)
	}
	if fromlen != 0 {
		if err := proc.Aspace.Userwriten(socklenn, 8, fromlen); err != 0 {
			return int(err)
		}
	}
	newfd := &vm.Fd_t{Fops: newfops}
	ret, ok := proc.Fd_insert(newfd, vm.FD_READ|vm.FD_WRITE)
	if !ok {
		vm.Close_panic(newfd)
		return int(-defs.EMFILE)
	}
	return ret
}

func copysockaddr(proc *common.Proc_t, san, sl int) ([]uint8, defs.Err_t) {
	if sl == 0 {
		return nil, 0
	}
	if sl < 0 {
		return nil, -defs.EFAULT
	}
	maxsl := 256
	if sl >= maxsl {
		return nil, -defs.ENOTSOCK
	}
	ub := proc.Aspace.Mkuserbuf(san, sl)
	sabuf := make([]uint8, sl)
	_, err := ub.Uioread(sabuf)
	//common.Ubpool.Put(ub)
	if err != 0 {
		return nil, err
	}
	return sabuf, 0
}

func sys_sendto(proc *common.Proc_t, fdn, bufn, flaglen, sockaddrn, socklen int) int {
	fd, err := _fd_write(proc, fdn)
	if err != 0 {
		return int(err)
	}
	flags := int(uint(uint32(flaglen)))
	if flags != 0 {
		panic("no imp")
	}
	buflen := int(uint(flaglen) >> 32)
	if buflen < 0 {
		return int(-defs.EFAULT)
	}

	// copy sockaddr to kernel space to avoid races
	sabuf, err := copysockaddr(proc, sockaddrn, socklen)
	if err != 0 {
		return int(err)
	}

	buf := proc.Aspace.Mkuserbuf(bufn, buflen)
	ret, err := fd.Fops.Sendmsg(buf, sabuf, nil, flags)
	//common.Ubpool.Put(buf)
	if err != 0 {
		return int(err)
	}
	return ret
}

func sys_recvfrom(proc *common.Proc_t, fdn, bufn, flaglen, sockaddrn,
	socklenn int) int {
	fd, err := _fd_read(proc, fdn)
	if err != 0 {
		return int(err)
	}
	flags := uint(uint32(flaglen))
	if flags != 0 {
		panic("no imp")
	}
	buflen := int(uint(flaglen) >> 32)
	buf := proc.Aspace.Mkuserbuf(bufn, buflen)

	// is the from address requested?
	var salen int
	if socklenn != 0 {
		l, err := proc.Aspace.Userreadn(socklenn, 8)
		if err != 0 {
			return int(err)
		}
		salen = l
		if salen < 0 {
			return int(-defs.EFAULT)
		}
	}
	fromsa := proc.Aspace.Mkuserbuf(sockaddrn, salen)
	ret, addrlen, _, _, err := fd.Fops.Recvmsg(buf, fromsa, zeroubuf, 0)
	//common.Ubpool.Put(buf)
	//common.Ubpool.Put(fromsa)
	if err != 0 {
		return int(err)
	}
	// write new socket size to user space
	if addrlen > 0 {
		if err := proc.Aspace.Userwriten(socklenn, 8, addrlen); err != 0 {
			return int(err)
		}
	}
	return ret
}

func sys_recvmsg(proc *common.Proc_t, fdn, _msgn, _flags int) int {
	if _flags != 0 {
		panic("no imp")
	}
	fd, err := _fd_read(proc, fdn)
	if err != 0 {
		return int(err)
	}
	// maybe copy the msghdr to kernel space?
	msgn := uint(_msgn)
	iovn, err1 := proc.Aspace.Userreadn(int(msgn+2*8), 8)
	niov, err2 := proc.Aspace.Userreadn(int(msgn+3*8), 4)
	cmsgl, err3 := proc.Aspace.Userreadn(int(msgn+5*8), 8)
	salen, err4 := proc.Aspace.Userreadn(int(msgn+1*8), 8)
	if err1 != 0 {
		return int(err1)
	}
	if err2 != 0 {
		return int(err2)
	}
	if err3 != 0 {
		return int(err3)
	}
	if err4 != 0 {
		return int(err4)
	}

	var saddr vm.Userio_i
	saddr = zeroubuf
	if salen > 0 {
		saddrn, err := proc.Aspace.Userreadn(int(msgn+0*8), 8)
		if err != 0 {
			return int(err)
		}
		ub := proc.Aspace.Mkuserbuf(saddrn, salen)
		saddr = ub
	}
	var cmsg vm.Userio_i
	cmsg = zeroubuf
	if cmsgl > 0 {
		cmsgn, err := proc.Aspace.Userreadn(int(msgn+4*8), 8)
		if err != 0 {
			return int(err)
		}
		ub := proc.Aspace.Mkuserbuf(cmsgn, cmsgl)
		cmsg = ub
	}

	iov := &vm.Useriovec_t{}
	err = iov.Iov_init(&proc.Aspace, uint(iovn), niov)
	if err != 0 {
		return int(err)
	}

	ret, sawr, cmwr, msgfl, err := fd.Fops.Recvmsg(iov, saddr,
		cmsg, 0)
	if err != 0 {
		return int(err)
	}
	// write size of socket address, ancillary data, and msg flags back to
	// user space
	if err := proc.Aspace.Userwriten(int(msgn+28), 4, int(msgfl)); err != 0 {
		return int(err)
	}
	if saddr.Totalsz() != 0 {
		if err := proc.Aspace.Userwriten(int(msgn+1*8), 8, sawr); err != 0 {
			return int(err)
		}
	}
	if cmsg.Totalsz() != 0 {
		if err := proc.Aspace.Userwriten(int(msgn+5*8), 8, cmwr); err != 0 {
			return int(err)
		}
	}
	return ret
}

func sys_sendmsg(proc *common.Proc_t, fdn, _msgn, _flags int) int {
	if _flags != 0 {
		panic("no imp")
	}
	fd, err := _fd_write(proc, fdn)
	if err != 0 {
		return int(err)
	}
	// maybe copy the msghdr to kernel space?
	msgn := uint(_msgn)
	iovn, err1 := proc.Aspace.Userreadn(int(msgn+2*8), 8)
	niov, err2 := proc.Aspace.Userreadn(int(msgn+3*8), 8)
	cmsgl, err3 := proc.Aspace.Userreadn(int(msgn+5*8), 8)
	salen, err4 := proc.Aspace.Userreadn(int(msgn+1*8), 8)
	if err1 != 0 {
		return int(err1)
	}
	if err2 != 0 {
		return int(err2)
	}
	if err3 != 0 {
		return int(err3)
	}
	if err4 != 0 {
		return int(err4)
	}

	// copy to address and ancillary data to kernel space
	var saddr []uint8
	if salen > 0 {
		if salen > 64 {
			return int(-defs.EINVAL)
		}
		saddrva, err := proc.Aspace.Userreadn(int(msgn+0*8), 8)
		if err != 0 {
			return int(err)
		}
		saddr = make([]uint8, salen)
		ub := proc.Aspace.Mkuserbuf(saddrva, salen)
		did, err := ub.Uioread(saddr)
		if err != 0 {
			return int(err)
		}
		if did != salen {
			panic("how")
		}
	}
	var cmsg []uint8
	if cmsgl > 0 {
		if cmsgl > 256 {
			return int(-defs.EINVAL)
		}
		cmsgva, err := proc.Aspace.Userreadn(int(msgn+4*8), 8)
		if err != 0 {
			return int(err)
		}
		cmsg = make([]uint8, cmsgl)
		ub := proc.Aspace.Mkuserbuf(cmsgva, cmsgl)
		did, err := ub.Uioread(cmsg)
		if err != 0 {
			return int(err)
		}
		if did != cmsgl {
			panic("how")
		}
	}
	iov := &vm.Useriovec_t{}
	err = iov.Iov_init(&proc.Aspace, uint(iovn), niov)
	if err != 0 {
		return int(err)
	}
	ret, err := fd.Fops.Sendmsg(iov, saddr, cmsg, 0)
	if err != 0 {
		return int(err)
	}
	return ret
}

func sys_socketpair(proc *common.Proc_t, domain, typ, proto int, sockn int) int {
	var opts common.Fdopt_t
	if typ&common.SOCK_NONBLOCK != 0 {
		opts |= common.O_NONBLOCK
	}
	var clop int
	if typ&common.SOCK_CLOEXEC != 0 {
		clop = vm.FD_CLOEXEC
	}

	mask := common.SOCK_STREAM | common.SOCK_DGRAM
	if typ&mask == 0 || typ&mask == mask {
		return int(-defs.EINVAL)
	}

	if !limits.Syslimit.Socks.Take() {
		return int(-defs.ENOMEM)
	}

	var sfops1, sfops2 vm.Fdops_i
	var err defs.Err_t
	switch {
	case domain == common.AF_UNIX && typ&common.SOCK_STREAM != 0:
		sfops1, sfops2, err = _suspair(opts)
	default:
		panic("no imp")
	}

	if err != 0 {
		limits.Syslimit.Socks.Give()
		return int(err)
	}

	fd1 := &vm.Fd_t{}
	fd1.Fops = sfops1
	fd2 := &vm.Fd_t{}
	fd2.Fops = sfops2
	perms := vm.FD_READ | vm.FD_WRITE | clop
	fdn1, fdn2, ok := proc.Fd_insert2(fd1, perms, fd2, perms)
	if !ok {
		vm.Close_panic(fd1)
		vm.Close_panic(fd2)
		return int(-defs.EMFILE)
	}
	if err1, err2 := proc.Aspace.Userwriten(sockn, 4, fdn1), proc.Aspace.Userwriten(sockn+4, 4, fdn2); err1 != 0 || err2 != 0 {
		if sys.Sys_close(proc, fdn1) != 0 || sys.Sys_close(proc, fdn2) != 0 {
			panic("must succeed")
		}
		if err1 == 0 {
			err1 = err2
		}
		return int(err1)
	}
	return 0
}

func _suspair(opts common.Fdopt_t) (vm.Fdops_i, vm.Fdops_i, defs.Err_t) {
	pipe1 := &pipe_t{}
	pipe2 := &pipe_t{}
	pipe1.pipe_start()
	pipe2.pipe_start()

	p1r := &pipefops_t{pipe: pipe1, options: opts}
	p1w := &pipefops_t{pipe: pipe2, writer: true, options: opts}

	p2r := &pipefops_t{pipe: pipe2, options: opts}
	p2w := &pipefops_t{pipe: pipe1, writer: true, options: opts}

	sfops1 := &susfops_t{pipein: p1r, pipeout: p1w, options: opts}
	sfops2 := &susfops_t{pipein: p2r, pipeout: p2w, options: opts}
	sfops1.conn, sfops2.conn = true, true
	return sfops1, sfops2, 0
}

func sys_shutdown(proc *common.Proc_t, fdn, how int) int {
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return int(-defs.EBADF)
	}
	var rdone, wdone bool
	if how&common.SHUT_WR != 0 {
		wdone = true
	}
	if how&common.SHUT_RD != 0 {
		rdone = true
	}
	return int(fd.Fops.Shutdown(rdone, wdone))
}

func sys_bind(proc *common.Proc_t, fdn, sockaddrn, socklen int) int {
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return int(-defs.EBADF)
	}

	sabuf, err := copysockaddr(proc, sockaddrn, socklen)
	if err != 0 {
		return int(err)
	}
	r := fd.Fops.Bind(sabuf)
	return int(r)
}

type sudfops_t struct {
	// this lock protects open and bound; bud has its own lock
	sync.Mutex
	bud   *bud_t
	open  int
	bound bool
}

func (sf *sudfops_t) Close() defs.Err_t {
	// XXX use new method
	sf.Lock()
	sf.open--
	if sf.open < 0 {
		panic("negative ref count")
	}
	term := sf.open == 0
	if term {
		if sf.bound {
			sf.bud.bud_close()
			sf.bound = false
			sf.bud = nil
		}
		limits.Syslimit.Socks.Give()
	}
	sf.Unlock()
	return 0
}

func (sf *sudfops_t) Fstat(s *stat.Stat_t) defs.Err_t {
	panic("no imp")
}

func (sf *sudfops_t) Mmapi(int, int, bool) ([]mem.Mmapinfo_t, defs.Err_t) {
	return nil, -defs.EINVAL
}

func (sf *sudfops_t) Pathi() defs.Inum_t {
	panic("cwd socket?")
}

func (sf *sudfops_t) Read(dst vm.Userio_i) (int, defs.Err_t) {
	return 0, -defs.EBADF
}

func (sf *sudfops_t) Reopen() defs.Err_t {
	sf.Lock()
	sf.open++
	sf.Unlock()
	return 0
}

func (sf *sudfops_t) Write(vm.Userio_i) (int, defs.Err_t) {
	return 0, -defs.EBADF
}

func (sf *sudfops_t) Truncate(newlen uint) defs.Err_t {
	return -defs.EINVAL
}

func (sf *sudfops_t) Pread(dst vm.Userio_i, offset int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

func (sf *sudfops_t) Pwrite(src vm.Userio_i, offset int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

func (sf *sudfops_t) Lseek(int, int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

// trims trailing nulls from slice
func slicetostr(buf []uint8) string {
	end := 0
	for i := range buf {
		end = i
		if buf[i] == 0 {
			break
		}
	}
	return string(buf[:end])
}

func (sf *sudfops_t) Accept(vm.Userio_i) (vm.Fdops_i, int, defs.Err_t) {
	return nil, 0, -defs.EINVAL
}

func (sf *sudfops_t) Bind(sa []uint8) defs.Err_t {
	sf.Lock()
	defer sf.Unlock()

	if sf.bound {
		return -defs.EINVAL
	}

	poff := 2
	path := ustr.MkUstrSlice(sa[poff:])
	// try to create the specified file as a special device
	bid := allbuds.bud_id_new()
	fsf, err := thefs.Fs_open_inner(path, common.O_CREAT|common.O_EXCL, 0, common.CurrentProc().Cwd, defs.D_SUD, int(bid))
	if err != 0 {
		return err
	}
	inum := fsf.Inum
	bud := allbuds.bud_new(bid, path, inum)
	if thefs.Fs_close(fsf.Inum) != 0 {
		panic("must succeed")
	}
	sf.bud = bud
	sf.bound = true
	return 0
}

func (sf *sudfops_t) Connect(sabuf []uint8) defs.Err_t {
	return -defs.EINVAL
}

func (sf *sudfops_t) Listen(backlog int) (vm.Fdops_i, defs.Err_t) {
	return nil, -defs.EINVAL
}

func (sf *sudfops_t) Sendmsg(src vm.Userio_i, sa []uint8,
	cmsg []uint8, flags int) (int, defs.Err_t) {
	if len(cmsg) != 0 || flags != 0 {
		panic("no imp")
	}
	poff := 2
	if len(sa) <= poff {
		return 0, -defs.EINVAL
	}
	st := &stat.Stat_t{}
	path := ustr.MkUstrSlice(sa[poff:])

	err := thefs.Fs_stat(path, st, common.CurrentProc().Cwd)
	if err != 0 {
		return 0, err
	}
	maj, min := unmkdev(st.Rdev())
	if maj != defs.D_SUD {
		return 0, -defs.ECONNREFUSED
	}
	ino := st.Rino()

	bid := budid_t(min)
	bud, ok := allbuds.bud_lookup(bid, defs.Inum_t(ino))
	if !ok {
		return 0, -defs.ECONNREFUSED
	}

	var bp ustr.Ustr
	sf.Lock()
	if sf.bound {
		bp = sf.bud.bpath
	}
	sf.Unlock()

	did, err := bud.bud_in(src, bp, cmsg)
	if err != 0 {
		return 0, err
	}
	return did, 0
}

func (sf *sudfops_t) Recvmsg(dst vm.Userio_i,
	fromsa vm.Userio_i, cmsg vm.Userio_i, flags int) (int, int, int, defs.Msgfl_t, defs.Err_t) {
	if cmsg.Totalsz() != 0 || flags != 0 {
		panic("no imp")
	}

	sf.Lock()
	defer sf.Unlock()

	// XXX what is recv'ing on an unbound unix datagram socket supposed to
	// do? openbsd and linux seem to block forever.
	if !sf.bound {
		return 0, 0, 0, 0, -defs.ECONNREFUSED
	}
	bud := sf.bud

	datadid, addrdid, ancdid, msgfl, err := bud.bud_out(dst, fromsa, cmsg)
	if err != 0 {
		return 0, 0, 0, 0, err
	}
	return datadid, addrdid, ancdid, msgfl, 0
}

func (sf *sudfops_t) Pollone(pm vm.Pollmsg_t) (vm.Ready_t, defs.Err_t) {
	sf.Lock()
	defer sf.Unlock()

	if !sf.bound {
		return pm.Events & vm.R_ERROR, 0
	}
	r, err := sf.bud.bud_poll(pm)
	return r, err
}

func (sf *sudfops_t) Fcntl(cmd, opt int) int {
	return int(-defs.ENOSYS)
}

func (sf *sudfops_t) Getsockopt(opt int, bufarg vm.Userio_i,
	intarg int) (int, defs.Err_t) {
	return 0, -defs.EOPNOTSUPP
}

func (sf *sudfops_t) Setsockopt(int, int, vm.Userio_i, int) defs.Err_t {
	return -defs.EOPNOTSUPP
}

func (sf *sudfops_t) Shutdown(read, write bool) defs.Err_t {
	return -defs.ENOTSOCK
}

type budid_t int

var allbuds = allbud_t{m: make(map[budkey_t]*bud_t)}

// buds are indexed by bid and inode number in order to detect stale socket
// files that happen to have the same bid.
type budkey_t struct {
	bid  budid_t
	priv defs.Inum_t
}

type allbud_t struct {
	// leaf lock
	sync.Mutex
	m       map[budkey_t]*bud_t
	nextbid budid_t
}

func (ab *allbud_t) bud_lookup(bid budid_t, fpriv defs.Inum_t) (*bud_t, bool) {
	key := budkey_t{bid, fpriv}

	ab.Lock()
	bud, ok := ab.m[key]
	ab.Unlock()

	return bud, ok
}

func (ab *allbud_t) bud_id_new() budid_t {
	ab.Lock()
	ret := ab.nextbid
	ab.nextbid++
	ab.Unlock()
	return ret
}

func (ab *allbud_t) bud_new(bid budid_t, budpath ustr.Ustr, fpriv defs.Inum_t) *bud_t {
	ret := &bud_t{}
	ret.bud_init(bid, budpath, fpriv)

	key := budkey_t{bid, fpriv}
	ab.Lock()
	if _, ok := ab.m[key]; ok {
		panic("bud exists")
	}
	ab.m[key] = ret
	ab.Unlock()
	return ret
}

func (ab *allbud_t) bud_del(bid budid_t, fpriv defs.Inum_t) {
	key := budkey_t{bid, fpriv}
	ab.Lock()
	if _, ok := ab.m[key]; !ok {
		panic("no such bud")
	}
	delete(ab.m, key)
	ab.Unlock()
}

type dgram_t struct {
	from ustr.Ustr
	sz   int
}

// a circular buffer for datagrams and their source addresses
type dgrambuf_t struct {
	cbuf   circbuf_t
	dgrams []dgram_t
	// add dgrams at head, remove from tail
	head uint
	tail uint
}

func (db *dgrambuf_t) dg_init(sz int) {
	db.cbuf.cb_init(sz)
	// assume that messages are at least 10 bytes
	db.dgrams = make([]dgram_t, sz/10)
	db.head, db.tail = 0, 0
}

// returns true if there is enough buffers to hold a datagram of size sz
func (db *dgrambuf_t) _canhold(sz int) bool {
	if (db.head-db.tail) == uint(len(db.dgrams)) ||
		db.cbuf.left() < sz {
		return false
	}
	return true
}

func (db *dgrambuf_t) _havedgram() bool {
	return db.head != db.tail
}

func (db *dgrambuf_t) copyin(src vm.Userio_i, from ustr.Ustr) (int, defs.Err_t) {
	// is there a free source address slot and buffer space?
	if !db._canhold(src.Totalsz()) {
		panic("should have blocked")
	}
	did, err := db.cbuf.copyin(src)
	if err != 0 {
		return 0, err
	}
	slot := &db.dgrams[db.head%uint(len(db.dgrams))]
	db.head++
	slot.from = from
	slot.sz = did
	return did, 0
}

func (db *dgrambuf_t) copyout(dst, fromsa, cmsg vm.Userio_i) (int, int, defs.Err_t) {
	if cmsg.Totalsz() != 0 {
		panic("no imp")
	}
	if db.head == db.tail {
		panic("should have blocked")
	}
	slot := &db.dgrams[db.tail%uint(len(db.dgrams))]
	sz := slot.sz
	if sz == 0 {
		panic("huh?")
	}
	var fdid int
	if fromsa.Totalsz() != 0 {
		fsaddr := _sockaddr_un(slot.from)
		var err defs.Err_t
		fdid, err = fromsa.Uiowrite(fsaddr)
		if err != 0 {
			return 0, 0, err
		}
	}
	did, err := db.cbuf.copyout_n(dst, sz)
	if err != 0 {
		return 0, 0, err
	}
	// commit tail
	db.tail++
	return did, fdid, 0
}

func (db *dgrambuf_t) dg_release() {
	db.cbuf.cb_release()
}

// convert bound socket path to struct sockaddr_un
func _sockaddr_un(budpath ustr.Ustr) []uint8 {
	ret := make([]uint8, 2, 16)
	// len
	writen(ret, 1, 0, len(budpath))
	// family
	writen(ret, 1, 1, common.AF_UNIX)
	// path
	ret = append(ret, budpath...)
	ret = append(ret, 0)
	return ret
}

// a type for bound UNIX datagram sockets
type bud_t struct {
	sync.Mutex
	bid     budid_t
	fpriv   defs.Inum_t
	dbuf    dgrambuf_t
	pollers vm.Pollers_t
	cond    *sync.Cond
	closed  bool
	bpath   ustr.Ustr
}

func (bud *bud_t) bud_init(bid budid_t, bpath ustr.Ustr, priv defs.Inum_t) {
	bud.bid = bid
	bud.fpriv = priv
	bud.bpath = bpath
	bud.dbuf.dg_init(512)
	bud.cond = sync.NewCond(bud)
}

func (bud *bud_t) _rready() {
	bud.cond.Broadcast()
	bud.pollers.Wakeready(vm.R_READ)
}

func (bud *bud_t) _wready() {
	bud.cond.Broadcast()
	bud.pollers.Wakeready(vm.R_WRITE)
}

// returns number of bytes written and error
func (bud *bud_t) bud_in(src vm.Userio_i, from ustr.Ustr, cmsg []uint8) (int, defs.Err_t) {
	if len(cmsg) != 0 {
		panic("no imp")
	}
	need := src.Totalsz()
	bud.Lock()
	for {
		if bud.closed {
			bud.Unlock()
			return 0, -defs.EBADF
		}
		if bud.dbuf._canhold(need) || bud.closed {
			break
		}
		if err := common.KillableWait(bud.cond); err != 0 {
			bud.Unlock()
			return 0, err
		}
	}
	did, err := bud.dbuf.copyin(src, from)
	bud._rready()
	bud.Unlock()
	return did, err
}

// returns number of bytes written of data, socket address, ancillary data, and
// ancillary message flags...
func (bud *bud_t) bud_out(dst, fromsa, cmsg vm.Userio_i) (int, int, int,
	defs.Msgfl_t, defs.Err_t) {
	if cmsg.Totalsz() != 0 {
		panic("no imp")
	}
	bud.Lock()
	for {
		if bud.closed {
			bud.Unlock()
			return 0, 0, 0, 0, -defs.EBADF
		}
		if bud.dbuf._havedgram() {
			break
		}
		if err := common.KillableWait(bud.cond); err != 0 {
			bud.Unlock()
			return 0, 0, 0, 0, err
		}
	}
	ddid, fdid, err := bud.dbuf.copyout(dst, fromsa, cmsg)
	bud._wready()
	bud.Unlock()
	return ddid, fdid, 0, 0, err
}

func (bud *bud_t) bud_poll(pm vm.Pollmsg_t) (vm.Ready_t, defs.Err_t) {
	var ret vm.Ready_t
	var err defs.Err_t
	bud.Lock()
	if bud.closed {
		goto out
	}
	if pm.Events&vm.R_READ != 0 && bud.dbuf._havedgram() {
		ret |= vm.R_READ
	}
	if pm.Events&vm.R_WRITE != 0 && bud.dbuf._canhold(32) {
		ret |= vm.R_WRITE
	}
	if ret == 0 && pm.Dowait {
		err = bud.pollers.Addpoller(&pm)
	}
out:
	bud.Unlock()
	return ret, err
}

// the bud is closed; wake up any waiting threads
func (bud *bud_t) bud_close() {
	bud.Lock()
	bud.closed = true
	bud.cond.Broadcast()
	bud.pollers.Wakeready(vm.R_READ | vm.R_WRITE | vm.R_ERROR)
	bid := bud.bid
	fpriv := bud.fpriv
	bud.dbuf.dg_release()
	bud.Unlock()

	allbuds.bud_del(bid, fpriv)
}

type susfops_t struct {
	pipein  *pipefops_t
	pipeout *pipefops_t
	bl      sync.Mutex
	conn    bool
	bound   bool
	lstn    bool
	myaddr  ustr.Ustr
	mysid   int
	options common.Fdopt_t
}

func (sus *susfops_t) Close() defs.Err_t {
	if !sus.conn {
		return 0
	}
	err1 := sus.pipein.Close()
	err2 := sus.pipeout.Close()
	if err1 != 0 {
		return err1
	}
	// XXX
	sus.pipein.pipe.Lock()
	term := sus.pipein.pipe.closed
	sus.pipein.pipe.Unlock()
	if term {
		limits.Syslimit.Socks.Give()
	}
	return err2
}

func (sus *susfops_t) Fstat(*stat.Stat_t) defs.Err_t {
	panic("no imp")
}

func (sus *susfops_t) Lseek(int, int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

func (sus *susfops_t) Mmapi(int, int, bool) ([]mem.Mmapinfo_t, defs.Err_t) {
	return nil, -defs.ENODEV
}

func (sus *susfops_t) Pathi() defs.Inum_t {
	panic("unix stream cwd?")
}

func (sus *susfops_t) Read(dst vm.Userio_i) (int, defs.Err_t) {
	read, _, _, _, err := sus.Recvmsg(dst, zeroubuf, zeroubuf, 0)
	return read, err
}

func (sus *susfops_t) Reopen() defs.Err_t {
	if !sus.conn {
		return 0
	}
	err1 := sus.pipein.Reopen()
	err2 := sus.pipeout.Reopen()
	if err1 != 0 {
		return err1
	}
	return err2
}

func (sus *susfops_t) Write(src vm.Userio_i) (int, defs.Err_t) {
	wrote, err := sus.Sendmsg(src, nil, nil, 0)
	if err == -defs.EPIPE {
		err = -defs.ECONNRESET
	}
	return wrote, err
}

func (sus *susfops_t) Truncate(newlen uint) defs.Err_t {
	return -defs.EINVAL
}

func (sus *susfops_t) Pread(dst vm.Userio_i, offset int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

func (sus *susfops_t) Pwrite(src vm.Userio_i, offset int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

func (sus *susfops_t) Accept(vm.Userio_i) (vm.Fdops_i, int, defs.Err_t) {
	return nil, 0, -defs.EINVAL
}

func (sus *susfops_t) Bind(saddr []uint8) defs.Err_t {
	sus.bl.Lock()
	defer sus.bl.Unlock()

	if sus.bound {
		return -defs.EINVAL
	}
	poff := 2
	path := ustr.MkUstrSlice(saddr[poff:])
	sid := susid_new()

	// create special file
	fsf, err := thefs.Fs_open_inner(path, common.O_CREAT|common.O_EXCL, 0, common.CurrentProc().Cwd, defs.D_SUS, sid)
	if err != 0 {
		return err
	}
	if thefs.Fs_close(fsf.Inum) != 0 {
		panic("must succeed")
	}
	sus.myaddr = path
	sus.mysid = sid
	sus.bound = true
	return 0
}

func (sus *susfops_t) Connect(saddr []uint8) defs.Err_t {
	sus.bl.Lock()
	defer sus.bl.Unlock()

	if sus.conn {
		return -defs.EISCONN
	}
	poff := 2
	path := ustr.MkUstrSlice(saddr[poff:])

	// lookup sid
	st := &stat.Stat_t{}
	err := thefs.Fs_stat(path, st, common.CurrentProc().Cwd)
	if err != 0 {
		return err
	}
	maj, min := unmkdev(st.Rdev())
	if maj != defs.D_SUS {
		return -defs.ECONNREFUSED
	}
	sid := min

	allsusl.Lock()
	susl, ok := allsusl.m[sid]
	allsusl.Unlock()
	if !ok {
		return -defs.ECONNREFUSED
	}

	pipein := &pipe_t{}
	pipein.pipe_start()

	pipeout, err := susl.connectwait(pipein)
	if err != 0 {
		return err
	}

	sus.pipein = &pipefops_t{pipe: pipein, options: sus.options}
	sus.pipeout = &pipefops_t{pipe: pipeout, writer: true, options: sus.options}
	sus.conn = true
	return 0
}

func (sus *susfops_t) Listen(backlog int) (vm.Fdops_i, defs.Err_t) {
	sus.bl.Lock()
	defer sus.bl.Unlock()

	if sus.conn {
		return nil, -defs.EISCONN
	}
	if !sus.bound {
		return nil, -defs.EINVAL
	}
	if sus.lstn {
		return nil, -defs.EINVAL
	}
	sus.lstn = true

	// create a listening socket
	susl := &susl_t{}
	susl.susl_start(sus.mysid, backlog)
	newsock := &suslfops_t{susl: susl, myaddr: sus.myaddr,
		options: sus.options}
	allsusl.Lock()
	// XXXPANIC
	if _, ok := allsusl.m[sus.mysid]; ok {
		panic("susl exists")
	}
	allsusl.m[sus.mysid] = susl
	allsusl.Unlock()

	return newsock, 0
}

func (sus *susfops_t) Sendmsg(src vm.Userio_i, toaddr []uint8,
	cmsg []uint8, flags int) (int, defs.Err_t) {
	if !sus.conn {
		return 0, -defs.ENOTCONN
	}
	if toaddr != nil {
		return 0, -defs.EISCONN
	}

	if len(cmsg) > 0 {
		scmsz := 16 + 8
		if len(cmsg) < scmsz {
			return 0, -defs.EINVAL
		}
		// allow fd sending
		cmsg_len := readn(cmsg, 8, 0)
		cmsg_level := readn(cmsg, 4, 8)
		cmsg_type := readn(cmsg, 4, 12)
		scm_rights := 1
		if cmsg_len != scmsz || cmsg_level != scm_rights ||
			cmsg_type != common.SOL_SOCKET {
			return 0, -defs.EINVAL
		}
		chdrsz := 16
		fdn := readn(cmsg, 4, chdrsz)
		ofd, ok := common.CurrentProc().Fd_get(fdn)
		if !ok {
			return 0, -defs.EBADF
		}
		nfd, err := vm.Copyfd(ofd)
		if err != 0 {
			return 0, err
		}
		err = sus.pipeout.pipe.op_fdadd(nfd)
		if err != 0 {
			return 0, err
		}
	}

	return sus.pipeout.Write(src)
}

func (sus *susfops_t) _fdrecv(cmsg vm.Userio_i,
	fl defs.Msgfl_t) (int, defs.Msgfl_t, defs.Err_t) {
	scmsz := 16 + 8
	if cmsg.Totalsz() < scmsz {
		return 0, fl, 0
	}
	nfd, ok := sus.pipein.pipe.op_fdtake()
	if !ok {
		return 0, fl, 0
	}
	nfdn, ok := common.CurrentProc().Fd_insert(nfd, nfd.Perms)
	if !ok {
		vm.Close_panic(nfd)
		return 0, fl, -defs.EMFILE
	}
	buf := make([]uint8, scmsz)
	writen(buf, 8, 0, scmsz)
	writen(buf, 4, 8, common.SOL_SOCKET)
	scm_rights := 1
	writen(buf, 4, 12, scm_rights)
	writen(buf, 4, 16, nfdn)
	l, err := cmsg.Uiowrite(buf)
	if err != 0 {
		return 0, fl, err
	}
	if l != scmsz {
		panic("how")
	}
	return scmsz, fl, 0
}

func (sus *susfops_t) Recvmsg(dst vm.Userio_i, fromsa vm.Userio_i,
	cmsg vm.Userio_i, flags int) (int, int, int, defs.Msgfl_t, defs.Err_t) {
	if !sus.conn {
		return 0, 0, 0, 0, -defs.ENOTCONN
	}

	ret, err := sus.pipein.Read(dst)
	if err != 0 {
		return 0, 0, 0, 0, err
	}
	cmsglen, msgfl, err := sus._fdrecv(cmsg, 0)
	return ret, 0, cmsglen, msgfl, err
}

func (sus *susfops_t) Pollone(pm vm.Pollmsg_t) (vm.Ready_t, defs.Err_t) {
	if !sus.conn {
		return pm.Events & vm.R_ERROR, 0
	}

	// pipefops_t.pollone() doesn't allow polling for reading on write-end
	// of pipe and vice versa
	var readyin vm.Ready_t
	var readyout vm.Ready_t
	both := pm.Events&(vm.R_READ|vm.R_WRITE) == 0
	var err defs.Err_t
	if both || pm.Events&vm.R_READ != 0 {
		readyin, err = sus.pipein.Pollone(pm)
	}
	if err != 0 {
		return 0, err
	}
	if readyin != 0 {
		return readyin, 0
	}
	if both || pm.Events&vm.R_WRITE != 0 {
		readyout, err = sus.pipeout.Pollone(pm)
	}
	return readyin | readyout, err
}

func (sus *susfops_t) Fcntl(cmd, opt int) int {
	sus.bl.Lock()
	defer sus.bl.Unlock()

	switch cmd {
	case common.F_GETFL:
		return int(sus.options)
	case common.F_SETFL:
		sus.options = common.Fdopt_t(opt)
		if sus.conn {
			sus.pipein.options = common.Fdopt_t(opt)
			sus.pipeout.options = common.Fdopt_t(opt)
		}
		return 0
	default:
		panic("weird cmd")
	}
}

func (sus *susfops_t) Getsockopt(opt int, bufarg vm.Userio_i,
	intarg int) (int, defs.Err_t) {
	switch opt {
	case common.SO_ERROR:
		dur := [4]uint8{}
		writen(dur[:], 4, 0, 0)
		did, err := bufarg.Uiowrite(dur[:])
		return did, err
	default:
		return 0, -defs.EOPNOTSUPP
	}
}

func (sus *susfops_t) Setsockopt(int, int, vm.Userio_i, int) defs.Err_t {
	return -defs.EOPNOTSUPP
}

func (sus *susfops_t) Shutdown(read, write bool) defs.Err_t {
	panic("no imp")
}

var _susid uint64

func susid_new() int {
	newid := atomic.AddUint64(&_susid, 1)
	return int(newid)
}

type allsusl_t struct {
	m map[int]*susl_t
	sync.Mutex
}

var allsusl = allsusl_t{m: map[int]*susl_t{}}

// listening unix stream socket
type susl_t struct {
	sync.Mutex
	waiters         []_suslblog_t
	pollers         vm.Pollers_t
	opencount       int
	mysid           int
	readyconnectors int
}

type _suslblog_t struct {
	conn *pipe_t
	acc  *pipe_t
	cond *sync.Cond
	err  defs.Err_t
}

func (susl *susl_t) susl_start(mysid, backlog int) {
	blm := 64
	if backlog < 0 || backlog > blm {
		backlog = blm
	}
	susl.waiters = make([]_suslblog_t, backlog)
	for i := range susl.waiters {
		susl.waiters[i].cond = sync.NewCond(susl)
	}
	susl.opencount = 1
	susl.mysid = mysid
}

func (susl *susl_t) _findbed(amconnector bool) (*_suslblog_t, bool) {
	for i := range susl.waiters {
		var chk *pipe_t
		if amconnector {
			chk = susl.waiters[i].conn
		} else {
			chk = susl.waiters[i].acc
		}
		if chk == nil {
			return &susl.waiters[i], true
		}
	}
	return nil, false
}

func (susl *susl_t) _findwaiter(getacceptor bool) (*_suslblog_t, bool) {
	for i := range susl.waiters {
		var chk *pipe_t
		var oth *pipe_t
		if getacceptor {
			chk = susl.waiters[i].acc
			oth = susl.waiters[i].conn
		} else {
			chk = susl.waiters[i].conn
			oth = susl.waiters[i].acc
		}
		if chk != nil && oth == nil {
			return &susl.waiters[i], true
		}
	}
	return nil, false
}

func (susl *susl_t) _slotreset(slot *_suslblog_t) {
	slot.acc = nil
	slot.conn = nil
}

func (susl *susl_t) _getpartner(mypipe *pipe_t, getacceptor,
	noblk bool) (*pipe_t, defs.Err_t) {
	susl.Lock()
	if susl.opencount == 0 {
		susl.Unlock()
		return nil, -defs.EBADF
	}

	var theirs *pipe_t
	// fastpath: is there a peer already waiting?
	s, found := susl._findwaiter(getacceptor)
	if found {
		if getacceptor {
			theirs = s.acc
			s.conn = mypipe
		} else {
			theirs = s.conn
			s.acc = mypipe
		}
		susl.Unlock()
		s.cond.Signal()
		return theirs, 0
	}
	if noblk {
		susl.Unlock()
		return nil, -defs.EWOULDBLOCK
	}
	// darn. wait for a peer.
	b, found := susl._findbed(getacceptor)
	if !found {
		// backlog is full
		susl.Unlock()
		if !getacceptor {
			panic("fixme: allow more accepts than backlog")
		}
		return nil, -defs.ECONNREFUSED
	}
	if getacceptor {
		b.conn = mypipe
		susl.pollers.Wakeready(vm.R_READ)
	} else {
		b.acc = mypipe
	}
	if getacceptor {
		susl.readyconnectors++
	}
	err := common.KillableWait(b.cond)
	if err == 0 {
		err = b.err
	}
	if getacceptor {
		theirs = b.acc
	} else {
		theirs = b.conn
	}
	susl._slotreset(b)
	if getacceptor {
		susl.readyconnectors--
	}
	susl.Unlock()
	return theirs, err
}

func (susl *susl_t) connectwait(mypipe *pipe_t) (*pipe_t, defs.Err_t) {
	noblk := false
	return susl._getpartner(mypipe, true, noblk)
}

func (susl *susl_t) acceptwait(mypipe *pipe_t) (*pipe_t, defs.Err_t) {
	noblk := false
	return susl._getpartner(mypipe, false, noblk)
}

func (susl *susl_t) acceptnowait(mypipe *pipe_t) (*pipe_t, defs.Err_t) {
	noblk := true
	return susl._getpartner(mypipe, false, noblk)
}

func (susl *susl_t) susl_reopen(delta int) defs.Err_t {
	ret := defs.Err_t(0)
	dorem := false
	susl.Lock()
	if susl.opencount != 0 {
		susl.opencount += delta
		if susl.opencount == 0 {
			dorem = true
		}
	} else {
		ret = -defs.EBADF
	}

	if dorem {
		limits.Syslimit.Socks.Give()
		// wake up all blocked connectors/acceptors/pollers
		for i := range susl.waiters {
			s := &susl.waiters[i]
			a := s.acc
			b := s.conn
			if a == nil && b == nil {
				continue
			}
			s.err = -defs.ECONNRESET
			s.cond.Signal()
		}
		susl.pollers.Wakeready(vm.R_READ | vm.R_HUP | vm.R_ERROR)
	}

	susl.Unlock()
	if dorem {
		allsusl.Lock()
		delete(allsusl.m, susl.mysid)
		allsusl.Unlock()
	}
	return ret
}

func (susl *susl_t) susl_poll(pm vm.Pollmsg_t) (vm.Ready_t, defs.Err_t) {
	susl.Lock()
	if susl.opencount == 0 {
		susl.Unlock()
		return 0, 0
	}
	if pm.Events&vm.R_READ != 0 {
		if susl.readyconnectors > 0 {
			susl.Unlock()
			return vm.R_READ, 0
		}
	}
	var err defs.Err_t
	if pm.Dowait {
		err = susl.pollers.Addpoller(&pm)
	}
	susl.Unlock()
	return 0, err
}

type suslfops_t struct {
	susl    *susl_t
	myaddr  ustr.Ustr
	options common.Fdopt_t
}

func (sf *suslfops_t) Close() defs.Err_t {
	return sf.susl.susl_reopen(-1)
}

func (sf *suslfops_t) Fstat(*stat.Stat_t) defs.Err_t {
	panic("no imp")
}

func (sf *suslfops_t) Lseek(int, int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

func (sf *suslfops_t) Mmapi(int, int, bool) ([]mem.Mmapinfo_t, defs.Err_t) {
	return nil, -defs.ENODEV
}

func (sf *suslfops_t) Pathi() defs.Inum_t {
	panic("unix stream listener cwd?")
}

func (sf *suslfops_t) Read(vm.Userio_i) (int, defs.Err_t) {
	return 0, -defs.ENOTCONN
}

func (sf *suslfops_t) Reopen() defs.Err_t {
	return sf.susl.susl_reopen(1)
}

func (sf *suslfops_t) Write(vm.Userio_i) (int, defs.Err_t) {
	return 0, -defs.EPIPE
}

func (sf *suslfops_t) Truncate(newlen uint) defs.Err_t {
	return -defs.EINVAL
}

func (sf *suslfops_t) Pread(dst vm.Userio_i, offset int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

func (sf *suslfops_t) Pwrite(src vm.Userio_i, offset int) (int, defs.Err_t) {
	return 0, -defs.ESPIPE
}

func (sf *suslfops_t) Accept(fromsa vm.Userio_i) (vm.Fdops_i, int, defs.Err_t) {
	// the connector has already taken syslimit.Socks (1 sock reservation
	// counts for a connected pair of UNIX stream sockets).
	noblk := sf.options&common.O_NONBLOCK != 0
	pipein := &pipe_t{}
	pipein.pipe_start()
	var pipeout *pipe_t
	var err defs.Err_t
	if noblk {
		pipeout, err = sf.susl.acceptnowait(pipein)
	} else {
		pipeout, err = sf.susl.acceptwait(pipein)
	}
	if err != 0 {
		return nil, 0, err
	}
	pfin := &pipefops_t{pipe: pipein, options: sf.options}
	pfout := &pipefops_t{pipe: pipeout, writer: true, options: sf.options}
	ret := &susfops_t{pipein: pfin, pipeout: pfout, conn: true,
		options: sf.options}
	return ret, 0, 0
}

func (sf *suslfops_t) Bind([]uint8) defs.Err_t {
	return -defs.EINVAL
}

func (sf *suslfops_t) Connect(sabuf []uint8) defs.Err_t {
	return -defs.EINVAL
}

func (sf *suslfops_t) Listen(backlog int) (vm.Fdops_i, defs.Err_t) {
	return nil, -defs.EINVAL
}

func (sf *suslfops_t) Sendmsg(vm.Userio_i, []uint8, []uint8,
	int) (int, defs.Err_t) {
	return 0, -defs.ENOTCONN
}

func (sf *suslfops_t) Recvmsg(vm.Userio_i, vm.Userio_i,
	vm.Userio_i, int) (int, int, int, defs.Msgfl_t, defs.Err_t) {
	return 0, 0, 0, 0, -defs.ENOTCONN
}

func (sf *suslfops_t) Pollone(pm vm.Pollmsg_t) (vm.Ready_t, defs.Err_t) {
	return sf.susl.susl_poll(pm)
}

func (sf *suslfops_t) Fcntl(cmd, opt int) int {
	switch cmd {
	case common.F_GETFL:
		return int(sf.options)
	case common.F_SETFL:
		sf.options = common.Fdopt_t(opt)
		return 0
	default:
		panic("weird cmd")
	}
}

func (sf *suslfops_t) Getsockopt(opt int, bufarg vm.Userio_i,
	intarg int) (int, defs.Err_t) {
	return 0, -defs.EOPNOTSUPP
}

func (sf *suslfops_t) Setsockopt(int, int, vm.Userio_i, int) defs.Err_t {
	return -defs.EOPNOTSUPP
}

func (sf *suslfops_t) Shutdown(read, write bool) defs.Err_t {
	return -defs.ENOTCONN
}

func sys_listen(proc *common.Proc_t, fdn, backlog int) int {
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return int(-defs.EBADF)
	}
	if backlog < 0 {
		backlog = 0
	}
	newfops, err := fd.Fops.Listen(backlog)
	if err != 0 {
		return int(err)
	}
	// replace old fops
	proc.Fdl.Lock()
	fd.Fops = newfops
	proc.Fdl.Unlock()
	return 0
}

func sys_getsockopt(proc *common.Proc_t, fdn, level, opt, optvaln, optlenn int) int {
	if level != common.SOL_SOCKET {
		panic("no imp")
	}
	var olen int
	if optlenn != 0 {
		l, err := proc.Aspace.Userreadn(optlenn, 8)
		if err != 0 {
			return int(err)
		}
		if l < 0 {
			return int(-defs.EFAULT)
		}
		olen = l
	}
	bufarg := proc.Aspace.Mkuserbuf(optvaln, olen)
	// XXX why intarg??
	intarg := optvaln
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return int(-defs.EBADF)
	}
	optwrote, err := fd.Fops.Getsockopt(opt, bufarg, intarg)
	if err != 0 {
		return int(err)
	}
	if optlenn != 0 {
		if err := proc.Aspace.Userwriten(optlenn, 8, optwrote); err != 0 {
			return int(err)
		}
	}
	return 0
}

func sys_setsockopt(proc *common.Proc_t, fdn, level, opt, optvaln, optlenn int) int {
	if optlenn < 0 {
		return int(-defs.EFAULT)
	}
	var intarg int
	if optlenn >= 4 {
		var err defs.Err_t
		intarg, err = proc.Aspace.Userreadn(optvaln, 4)
		if err != 0 {
			return int(err)
		}
	}
	bufarg := proc.Aspace.Mkuserbuf(optvaln, optlenn)
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return int(-defs.EBADF)
	}
	err := fd.Fops.Setsockopt(level, opt, bufarg, intarg)
	return int(err)
}

func _closefds(fds []*vm.Fd_t) {
	for _, fd := range fds {
		if fd != nil {
			vm.Close_panic(fd)
		}
	}
}

func sys_fork(parent *common.Proc_t, ptf *[common.TFSIZE]uintptr, tforkp int, flags int) int {
	tmp := flags & (common.FORK_THREAD | common.FORK_PROCESS)
	if tmp != common.FORK_THREAD && tmp != common.FORK_PROCESS {
		return int(-defs.EINVAL)
	}

	mkproc := flags&common.FORK_PROCESS != 0
	var child *common.Proc_t
	var childtid defs.Tid_t
	var ret int

	// copy parents trap frame
	chtf := &[common.TFSIZE]uintptr{}
	*chtf = *ptf

	if mkproc {
		var ok bool
		// lock fd table for copying
		parent.Fdl.Lock()
		cwd := *parent.Cwd
		child, ok = common.Proc_new(parent.Name, &cwd, parent.Fds, sys)
		parent.Fdl.Unlock()
		if !ok {
			lhits++
			return int(-defs.ENOMEM)
		}

		child.Aspace.Pmap, child.Aspace.P_pmap, ok = physmem.Pmap_new()
		if !ok {
			goto outproc
		}
		physmem.Refup(child.Aspace.P_pmap)

		child.Pwait = &parent.Mywait
		ok = parent.Start_proc(child.Pid)
		if !ok {
			lhits++
			goto outmem
		}

		// fork parent address space
		parent.Aspace.Lock_pmap()
		rsp := chtf[common.TF_RSP]
		doflush, ok := parent.Vm_fork(child, rsp)
		if ok && !doflush {
			panic("no writable segs?")
		}
		// flush all ptes now marked COW
		if doflush {
			parent.Tlbflush()
		}
		parent.Aspace.Unlock_pmap()

		if !ok {
			// child page table allocation failed. call
			// common.Proc_t.terminate which will clean everything up. the
			// parent will get th error code directly.
			child.Thread_dead(child.Tid0(), 0, false)
			return int(-defs.ENOMEM)
		}

		childtid = child.Tid0()
		ret = child.Pid
	} else {
		// validate tfork struct
		tcb, err1 := parent.Aspace.Userreadn(tforkp+0, 8)
		tidaddrn, err2 := parent.Aspace.Userreadn(tforkp+8, 8)
		stack, err3 := parent.Aspace.Userreadn(tforkp+16, 8)
		if err1 != 0 {
			return int(err1)
		}
		if err2 != 0 {
			return int(err2)
		}
		if err3 != 0 {
			return int(err3)
		}
		writetid := tidaddrn != 0
		if tcb != 0 {
			chtf[common.TF_FSBASE] = uintptr(tcb)
		}

		child = parent
		var ok bool
		childtid, ok = parent.Thread_new()
		if !ok {
			lhits++
			return int(-defs.ENOMEM)
		}
		ok = parent.Start_thread(childtid)
		if !ok {
			lhits++
			parent.Thread_undo(childtid)
			return int(-defs.ENOMEM)
		}

		v := int(childtid)
		chtf[common.TF_RSP] = uintptr(stack)
		ret = v
		if writetid {
			// it is not a fatal error if some thread unmapped the
			// memory that was supposed to hold the new thread's
			// tid out from under us.
			parent.Aspace.Userwriten(tidaddrn, 8, v)
		}
	}

	chtf[common.TF_RAX] = 0
	child.Sched_add(chtf, childtid)
	return ret
outmem:
	physmem.Refdown(child.Aspace.P_pmap)
outproc:
	common.Tid_del()
	common.Proc_del(child.Pid)
	_closefds(child.Fds)
	return int(-defs.ENOMEM)
}

func sys_execv(proc *common.Proc_t, tf *[common.TFSIZE]uintptr, pathn int, argn int) int {
	args, err := proc.Userargs(argn)
	if err != 0 {
		return int(err)
	}
	path, err := proc.Aspace.Userstr(pathn, fs.NAME_MAX)
	if err != 0 {
		return int(err)
	}
	err = badpath(path)
	if err != 0 {
		return int(err)
	}
	return sys_execv1(proc, tf, path, args)
}

var _zvmregion vm.Vmregion_t

func sys_execv1(proc *common.Proc_t, tf *[common.TFSIZE]uintptr, paths ustr.Ustr,
	args []ustr.Ustr) int {
	// XXX a multithreaded process that execs is broken; POSIX2008 says
	// that all threads should terminate before exec.
	if proc.Thread_count() > 1 {
		panic("fix exec with many threads")
	}

	proc.Aspace.Lock_pmap()
	defer proc.Aspace.Unlock_pmap()

	// save page trackers in case the exec fails
	ovmreg := proc.Aspace.Vmregion
	proc.Aspace.Vmregion = _zvmregion

	// create kernel page table
	opmap := proc.Aspace.Pmap
	op_pmap := proc.Aspace.P_pmap
	var ok bool
	proc.Aspace.Pmap, proc.Aspace.P_pmap, ok = physmem.Pmap_new()
	if !ok {
		proc.Aspace.Pmap, proc.Aspace.P_pmap = opmap, op_pmap
		return int(-defs.ENOMEM)
	}
	physmem.Refup(proc.Aspace.P_pmap)
	for _, e := range mem.Kents {
		proc.Aspace.Pmap[e.Pml4slot] = e.Entry
	}

	restore := func() {
		vm.Uvmfree_inner(proc.Aspace.Pmap, proc.Aspace.P_pmap, &proc.Aspace.Vmregion)
		physmem.Refdown(proc.Aspace.P_pmap)
		proc.Aspace.Vmregion.Clear()
		proc.Aspace.Pmap = opmap
		proc.Aspace.P_pmap = op_pmap
		proc.Aspace.Vmregion = ovmreg
	}

	// load binary image -- get first block of file
	file, err := thefs.Fs_open(paths, common.O_RDONLY, 0, proc.Cwd, 0, 0)
	if err != 0 {
		restore()
		return int(err)
	}
	defer vm.Close_panic(file)

	hdata := make([]uint8, 512)
	ub := &vm.Fakeubuf_t{}
	ub.Fake_init(hdata)
	ret, err := file.Fops.Read(ub)
	if err != 0 {
		restore()
		return int(err)
	}
	if ret < len(hdata) {
		hdata = hdata[0:ret]
	}

	// assume its always an elf, for now
	elfhdr := &elf_t{hdata}
	ok = elfhdr.sanity()
	if !ok {
		restore()
		return int(-defs.EPERM)
	}

	// elf_load() will create two copies of TLS section: one for the fresh
	// copy and one for thread 0
	freshtls, t0tls, tlssz, err := elfhdr.elf_load(proc, file)
	if err != 0 {
		restore()
		return int(err)
	}

	// map new stack
	numstkpages := 6
	// +1 for the guard page
	stksz := (numstkpages + 1) * mem.PGSIZE
	stackva := proc.Aspace.Unusedva_inner(0x0ff<<39, stksz)
	proc.Aspace.Vmadd_anon(stackva, mem.PGSIZE, 0)
	proc.Aspace.Vmadd_anon(stackva+mem.PGSIZE, stksz-mem.PGSIZE, vm.PTE_U|vm.PTE_W)
	stackva += stksz
	// eagerly map first two pages for stack
	stkeagermap := 2
	for i := 0; i < stkeagermap; i++ {
		p := uintptr(stackva - (i+1)*mem.PGSIZE)
		_, p_pg, ok := physmem.Refpg_new()
		if !ok {
			restore()
			return int(-defs.ENOMEM)
		}
		_, ok = proc.Aspace.Page_insert(int(p), p_pg, vm.PTE_W|vm.PTE_U, true)
		if !ok {
			restore()
			return int(-defs.ENOMEM)
		}
	}

	// XXX make insertargs not fail by using more than a page...
	argc, argv, err := insertargs(proc, args)
	if err != 0 {
		restore()
		return int(err)
	}

	// put special struct on stack: fresh tls start, tls len, and tls0
	// pointer
	words := 4
	buf := make([]uint8, words*8)
	writen(buf, 8, 0, freshtls)
	writen(buf, 8, 8, tlssz)
	writen(buf, 8, 16, t0tls)
	writen(buf, 8, 24, int(runtime.Pspercycle))
	bufdest := stackva - words*8
	tls0addr := bufdest + 2*8

	if err := proc.Aspace.K2user_inner(buf, bufdest); err != 0 {
		restore()
		return int(err)
	}

	// the exec must succeed now; free old pmap/mapped files
	if op_pmap != 0 {
		vm.Uvmfree_inner(opmap, op_pmap, &ovmreg)
		physmem.Dec_pmap(op_pmap)
	}
	ovmreg.Clear()

	// close fds marked with CLOEXEC
	for fdn, fd := range proc.Fds {
		if fd == nil {
			continue
		}
		if fd.Perms&vm.FD_CLOEXEC != 0 {
			if sys.Sys_close(proc, fdn) != 0 {
				panic("close")
			}
		}
	}

	// commit new image state
	tf[common.TF_RSP] = uintptr(bufdest)
	tf[common.TF_RIP] = uintptr(elfhdr.entry())
	tf[common.TF_RFLAGS] = uintptr(common.TF_FL_IF)
	ucseg := uintptr(5)
	udseg := uintptr(6)
	tf[common.TF_CS] = (ucseg << 3) | 3
	tf[common.TF_SS] = (udseg << 3) | 3
	tf[common.TF_RDI] = uintptr(argc)
	tf[common.TF_RSI] = uintptr(argv)
	tf[common.TF_RDX] = uintptr(bufdest)
	tf[common.TF_FSBASE] = uintptr(tls0addr)
	proc.Mmapi = mem.USERMIN
	proc.Name = paths

	return 0
}

func insertargs(proc *common.Proc_t, sargs []ustr.Ustr) (int, int, defs.Err_t) {
	// find free page
	uva := proc.Aspace.Unusedva_inner(0, mem.PGSIZE)
	proc.Aspace.Vmadd_anon(uva, mem.PGSIZE, vm.PTE_U)
	_, p_pg, ok := physmem.Refpg_new()
	if !ok {
		return 0, 0, -defs.ENOMEM
	}
	_, ok = proc.Aspace.Page_insert(uva, p_pg, vm.PTE_U, true)
	if !ok {
		physmem.Refdown(p_pg)
		return 0, 0, -defs.ENOMEM
	}
	//var args [][]uint8
	args := make([][]uint8, 0, 12)
	for _, str := range sargs {
		args = append(args, []uint8(str))
	}
	argptrs := make([]int, len(args)+1)
	// copy strings to arg page
	cnt := 0
	for i, arg := range args {
		argptrs[i] = uva + cnt
		// add null terminators
		arg = append(arg, 0)
		if err := proc.Aspace.K2user_inner(arg, uva+cnt); err != 0 {
			// args take up more than a page? the user is on their
			// own.
			return 0, 0, err
		}
		cnt += len(arg)
	}
	argptrs[len(argptrs)-1] = 0
	// now put the array of strings
	argstart := uva + cnt
	vdata, err := proc.Aspace.Userdmap8_inner(argstart, true)
	if err != 0 || len(vdata) < len(argptrs)*8 {
		fmt.Printf("no room for args")
		// XXX
		return 0, 0, -defs.ENOSPC
	}
	for i, ptr := range argptrs {
		writen(vdata, 8, i*8, ptr)
	}
	return len(args), argstart, 0
}

func (s *syscall_t) Sys_exit(proc *common.Proc_t, tid defs.Tid_t, status int) {
	// set doomed so all other threads die
	proc.Doomall()
	proc.Thread_dead(tid, status, true)
}

func sys_threxit(proc *common.Proc_t, tid defs.Tid_t, status int) {
	proc.Thread_dead(tid, status, false)
}

func sys_wait4(proc *common.Proc_t, tid defs.Tid_t, wpid, statusp, options, rusagep,
	_isthread int) int {
	if wpid == common.WAIT_MYPGRP || options == common.WCONTINUED ||
		options == common.WUNTRACED {
		panic("no imp")
	}

	// no waiting for yourself!
	if tid == defs.Tid_t(wpid) {
		return int(-defs.ECHILD)
	}
	isthread := _isthread != 0
	if isthread && wpid == common.WAIT_ANY {
		return int(-defs.EINVAL)
	}

	noblk := options&common.WNOHANG != 0
	var resp common.Waitst_t
	var err defs.Err_t
	if isthread {
		resp, err = proc.Mywait.Reaptid(wpid, noblk)
	} else {
		resp, err = proc.Mywait.Reappid(wpid, noblk)
	}

	if err != 0 {
		return int(err)
	}
	if isthread {
		if statusp != 0 {
			err := proc.Aspace.Userwriten(statusp, 8, resp.Status)
			if err != 0 {
				return int(err)
			}
		}
	} else {
		if statusp != 0 {
			err := proc.Aspace.Userwriten(statusp, 4, resp.Status)
			if err != 0 {
				return int(err)
			}
		}
		// update total child rusage
		proc.Catime.Add(&resp.Atime)
		if rusagep != 0 {
			ru := resp.Atime.To_rusage()
			err = proc.Aspace.K2user(ru, rusagep)
		}
		if err != 0 {
			return int(err)
		}
	}
	return resp.Pid
}

func sys_kill(proc *common.Proc_t, pid, sig int) int {
	if sig != common.SIGKILL {
		panic("no imp")
	}
	p, ok := common.Proc_check(pid)
	if !ok {
		return int(-defs.ESRCH)
	}
	p.Doomall()
	return 0
}

func sys_pread(proc *common.Proc_t, fdn, bufn, lenn, offset int) int {
	fd, err := _fd_read(proc, fdn)
	if err != 0 {
		return int(err)
	}
	dst := proc.Aspace.Mkuserbuf(bufn, lenn)
	ret, err := fd.Fops.Pread(dst, offset)
	if err != 0 {
		return int(err)
	}
	return ret
}

func sys_pwrite(proc *common.Proc_t, fdn, bufn, lenn, offset int) int {
	fd, err := _fd_write(proc, fdn)
	if err != 0 {
		return int(err)
	}
	src := proc.Aspace.Mkuserbuf(bufn, lenn)
	ret, err := fd.Fops.Pwrite(src, offset)
	if err != 0 {
		return int(err)
	}
	return ret
}

type futexmsg_t struct {
	op      uint
	aux     uint32
	ack     chan int
	othmut  futex_t
	cndtake []chan int
	totake  []_futto_t
	fumem   futumem_t
	timeout time.Time
	useto   bool
}

func (fm *futexmsg_t) fmsg_init(op uint, aux uint32, ack chan int) {
	fm.op = op
	fm.aux = aux
	fm.ack = ack
}

// futex timeout metadata
type _futto_t struct {
	when   time.Time
	tochan <-chan time.Time
	who    chan int
}

type futex_t struct {
	reopen chan int
	cmd    chan futexmsg_t
	_cnds  []chan int
	cnds   []chan int
	_tos   []_futto_t
	tos    []_futto_t
}

func (f *futex_t) cndsleep(c chan int) {
	f.cnds = append(f.cnds, c)
}

func (f *futex_t) cndwake(v int) {
	if len(f.cnds) == 0 {
		return
	}
	c := f.cnds[0]
	f.cnds = f.cnds[1:]
	if len(f.cnds) == 0 {
		f.cnds = f._cnds
	}
	f._torm(c)
	c <- v
}

func (f *futex_t) toadd(who chan int, when time.Time) {
	fto := _futto_t{when, time.After(when.Sub(time.Now())), who}
	f.tos = append(f.tos, fto)
}

func (f *futex_t) tonext() (<-chan time.Time, chan int) {
	if len(f.tos) == 0 {
		return nil, nil
	}
	small := f.tos[0].when
	next := f.tos[0]
	for _, nto := range f.tos {
		if nto.when.Before(small) {
			small = nto.when
			next = nto
		}
	}
	return next.tochan, next.who
}

func (f *futex_t) _torm(who chan int) {
	idx := -1
	for i, nto := range f.tos {
		if nto.who == who {
			idx = i
			break
		}
	}
	if idx == -1 {
		return
	}
	copy(f.tos[idx:], f.tos[idx+1:])
	l := len(f.tos)
	f.tos = f.tos[:l-1]
	if len(f.tos) == 0 {
		f.tos = f._tos
	}
}

func (f *futex_t) towake(who chan int, v int) {
	// remove from tos and cnds
	f._torm(who)
	idx := -1
	for i := range f.cnds {
		if f.cnds[i] == who {
			idx = i
			break
		}
	}
	copy(f.cnds[idx:], f.cnds[idx+1:])
	l := len(f.cnds)
	f.cnds = f.cnds[:l-1]
	if len(f.cnds) == 0 {
		f.cnds = f._cnds
	}
	who <- v
}

const (
	_FUTEX_LAST = common.FUTEX_CNDGIVE
	// futex internal op
	_FUTEX_CNDTAKE = 4
)

func (f *futex_t) _resume(ack chan int, err defs.Err_t) {
	select {
	case ack <- int(err):
	default:
	}
}

func (f *futex_t) futex_start() {
	res.Kresdebug(1<<10, "futex daemon")
	maxwait := 10
	f._cnds = make([]chan int, 0, maxwait)
	f.cnds = f._cnds
	f._tos = make([]_futto_t, 0, maxwait)
	f.tos = f._tos

	pack := make(chan int, 1)
	opencount := 1
	for opencount > 0 {
		res.Kunresdebug()
		res.Kresdebug(1<<10, "futex daemon")
		tochan, towho := f.tonext()
		select {
		case <-tochan:
			f.towake(towho, 0)
		case d := <-f.reopen:
			opencount += d
		case fm := <-f.cmd:
			switch fm.op {
			case common.FUTEX_SLEEP:
				val, err := fm.fumem.futload()
				if err != 0 {
					f._resume(fm.ack, err)
					break
				}
				if val != fm.aux {
					// owner just unlocked and it's this
					// thread's turn; don't sleep
					f._resume(fm.ack, 0)
				} else {
					if (fm.useto && len(f.tos) >= maxwait) ||
						len(f.cnds) >= maxwait {
						f._resume(fm.ack, -defs.ENOMEM)
						break
					}
					if fm.useto {
						f.toadd(fm.ack, fm.timeout)
					}
					f.cndsleep(fm.ack)
				}
			case common.FUTEX_WAKE:
				var v int
				if fm.aux == 1 {
					v = 0
				} else if fm.aux == ^uint32(0) {
					v = 1
				} else {
					panic("weird wake n")
				}
				f.cndwake(v)
				f._resume(fm.ack, 0)
			case common.FUTEX_CNDGIVE:
				// as an optimization to avoid thundering herd
				// after pthread_cond_broadcast(3), move
				// conditional variable's queue of sleepers to
				// the mutex of the thread we wakeup here.
				l := len(f.cnds)
				if l == 0 {
					f._resume(fm.ack, 0)
					break
				}
				here := make([]chan int, l)
				copy(here, f.cnds)
				tohere := make([]_futto_t, len(f.tos))
				copy(tohere, f.tos)

				var nfm futexmsg_t
				nfm.fmsg_init(_FUTEX_CNDTAKE, 0, pack)
				nfm.cndtake = here
				nfm.totake = tohere

				fm.othmut.cmd <- nfm
				err := <-nfm.ack
				if err == 0 {
					f.cnds = f._cnds
					f.tos = f._tos
				}
				f._resume(fm.ack, defs.Err_t(err))
			case _FUTEX_CNDTAKE:
				// add new waiters to our queue; get them
				// tickets
				here := fm.cndtake
				tohere := fm.totake
				if len(f.cnds)+len(here) >= maxwait ||
					len(f.tos)+len(tohere) >= maxwait {
					f._resume(fm.ack, -defs.ENOMEM)
					break
				}
				f.cnds = append(f.cnds, here...)
				f.tos = append(f.tos, tohere...)
				f._resume(fm.ack, 0)
			default:
				panic("bad futex op")
			}
		}
	}
	res.Kunresdebug()
}

type allfutex_t struct {
	sync.Mutex
	m map[uintptr]futex_t
}

var _allfutex = allfutex_t{m: map[uintptr]futex_t{}}

func futex_ensure(uniq uintptr) (futex_t, defs.Err_t) {
	_allfutex.Lock()
	if len(_allfutex.m) > limits.Syslimit.Futexes {
		_allfutex.Unlock()
		var zf futex_t
		return zf, -defs.ENOMEM
	}
	r, ok := _allfutex.m[uniq]
	if !ok {
		r.reopen = make(chan int)
		r.cmd = make(chan futexmsg_t)
		_allfutex.m[uniq] = r
		go r.futex_start()
	}
	_allfutex.Unlock()
	return r, 0
}

// pmap must be locked. maps user va to kernel va. returns kva as uintptr and
// *uint32
func _uva2kva(proc *common.Proc_t, va uintptr) (uintptr, *uint32, defs.Err_t) {
	proc.Aspace.Lockassert_pmap()

	pte := vm.Pmap_lookup(proc.Aspace.Pmap, int(va))
	if pte == nil || *pte&vm.PTE_P == 0 || *pte&vm.PTE_U == 0 {
		return 0, nil, -defs.EFAULT
	}
	pgva := physmem.Dmap(*pte & vm.PTE_ADDR)
	pgoff := uintptr(va) & uintptr(vm.PGOFFSET)
	uniq := uintptr(unsafe.Pointer(pgva)) + pgoff
	return uniq, (*uint32)(unsafe.Pointer(uniq)), 0
}

func va2fut(proc *common.Proc_t, va uintptr) (futex_t, defs.Err_t) {
	proc.Aspace.Lock_pmap()
	defer proc.Aspace.Unlock_pmap()

	var zf futex_t
	uniq, _, err := _uva2kva(proc, va)
	if err != 0 {
		return zf, err
	}
	return futex_ensure(uniq)
}

// an object for atomically looking-up and incrementing/loading from a user
// address
type futumem_t struct {
	proc *common.Proc_t
	umem uintptr
}

func (fu *futumem_t) futload() (uint32, defs.Err_t) {
	fu.proc.Aspace.Lock_pmap()
	defer fu.proc.Aspace.Unlock_pmap()

	_, ptr, err := _uva2kva(fu.proc, fu.umem)
	if err != 0 {
		return 0, err
	}
	var ret uint32
	ret = atomic.LoadUint32(ptr)
	return ret, 0
}

func sys_futex(proc *common.Proc_t, _op, _futn, _fut2n, aux, timespecn int) int {
	op := uint(_op)
	if op > _FUTEX_LAST {
		return int(-defs.EINVAL)
	}
	futn := uintptr(_futn)
	fut2n := uintptr(_fut2n)
	// futn must be 4 byte aligned
	if (futn|fut2n)&0x3 != 0 {
		return int(-defs.EINVAL)
	}
	fut, err := va2fut(proc, futn)
	if err != 0 {
		return int(err)
	}

	var fm futexmsg_t
	// could lazily allocate one futex channel per thread
	fm.fmsg_init(op, uint32(aux), make(chan int, 1))
	fm.fumem = futumem_t{proc, futn}

	if timespecn != 0 {
		_, when, err := proc.Aspace.Usertimespec(timespecn)
		if err != 0 {
			return int(err)
		}
		n := time.Now()
		if when.Before(n) {
			return int(-defs.EINVAL)
		}
		fm.timeout = when
		fm.useto = true
	}

	if op == common.FUTEX_CNDGIVE {
		fm.othmut, err = va2fut(proc, fut2n)
		if err != 0 {
			return int(err)
		}
	}

	kn := &tinfo.Current().Killnaps
	fut.cmd <- fm
	select {
	case ret := <-fm.ack:
		return ret
	case <-kn.Killch:
		if kn.Kerr == 0 {
			panic("no")
		}
		return int(kn.Kerr)
	}
}

func sys_gettid(proc *common.Proc_t, tid defs.Tid_t) int {
	return int(tid)
}

func sys_fcntl(proc *common.Proc_t, fdn, cmd, opt int) int {
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return int(-defs.EBADF)
	}
	switch cmd {
	// general fcntl(2) ops
	case common.F_GETFD:
		return fd.Perms & vm.FD_CLOEXEC
	case common.F_SETFD:
		if opt&vm.FD_CLOEXEC == 0 {
			fd.Perms &^= vm.FD_CLOEXEC
		} else {
			fd.Perms |= vm.FD_CLOEXEC
		}
		return 0
	// fd specific fcntl(2) ops
	case common.F_GETFL, common.F_SETFL:
		return fd.Fops.Fcntl(cmd, opt)
	default:
		return int(-defs.EINVAL)
	}
}

func sys_truncate(proc *common.Proc_t, pathn int, newlen uint) int {
	path, err := proc.Aspace.Userstr(pathn, fs.NAME_MAX)
	if err != 0 {
		return int(err)
	}
	if err := badpath(path); err != 0 {
		return int(err)
	}
	f, err := thefs.Fs_open(path, common.O_WRONLY, 0, proc.Cwd, 0, 0)
	if err != 0 {
		return int(err)
	}
	err = f.Fops.Truncate(newlen)
	vm.Close_panic(f)
	return int(err)
}

func sys_ftruncate(proc *common.Proc_t, fdn int, newlen uint) int {
	fd, ok := proc.Fd_get(fdn)
	if !ok {
		return int(-defs.EBADF)
	}
	return int(fd.Fops.Truncate(newlen))
}

func sys_getcwd(proc *common.Proc_t, bufn, sz int) int {
	dst := proc.Aspace.Mkuserbuf(bufn, sz)
	_, err := dst.Uiowrite([]uint8(proc.Cwd.Path))
	if err != 0 {
		return int(err)
	}
	if _, err := dst.Uiowrite([]uint8{0}); err != 0 {
		return int(err)
	}
	return 0
}

func sys_chdir(proc *common.Proc_t, dirn int) int {
	path, err := proc.Aspace.Userstr(dirn, fs.NAME_MAX)
	if err != 0 {
		return int(err)
	}
	err = badpath(path)
	if err != 0 {
		return int(err)
	}

	proc.Cwd.Lock()
	defer proc.Cwd.Unlock()

	newfd, err := thefs.Fs_open(path, common.O_RDONLY|common.O_DIRECTORY, 0, proc.Cwd, 0, 0)
	if err != 0 {
		return int(err)
	}
	vm.Close_panic(proc.Cwd.Fd)
	proc.Cwd.Fd = newfd
	if path.IsAbsolute() {
		proc.Cwd.Path = bpath.Canonicalize(path)
	} else {
		proc.Cwd.Path = proc.Cwd.Canonicalpath(path)
	}
	return 0
}

func badpath(path ustr.Ustr) defs.Err_t {
	if len(path) == 0 {
		return -defs.ENOENT
	}
	return 0
}

func buftodests(buf []uint8, dsts [][]uint8) int {
	ret := 0
	for _, dst := range dsts {
		ub := len(buf)
		if ub > len(dst) {
			ub = len(dst)
		}
		for i := 0; i < ub; i++ {
			dst[i] = buf[i]
		}
		ret += ub
		buf = buf[ub:]
	}
	return ret
}

func _prof_go(en bool) {
	if en {
		prof.init()
		err := pprof.StartCPUProfile(&prof)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		//runtime.SetBlockProfileRate(1)
	} else {
		pprof.StopCPUProfile()
		prof.dump()

		//pprof.WriteHeapProfile(&prof)
		//prof.dump()

		//p := pprof.Lookup("block")
		//err := p.WriteTo(&prof, 0)
		//if err != nil {
		//	fmt.Printf("%v\n", err)
		//	return
		//}
		//prof.dump()
	}
}

func _prof_nmi(en bool, pmev pmev_t, intperiod int) {
	if en {
		min := uint(intperiod)
		// default unhalted cycles sampling rate
		defperiod := intperiod == 0
		if defperiod && pmev.evid == EV_UNHALTED_CORE_CYCLES {
			cyc := runtime.Cpumhz * 1000000
			samples := uint(1000)
			min = cyc / samples
		}
		max := uint(float64(min) * 1.2)
		if !profhw.startnmi(pmev.evid, pmev.pflags, min, max) {
			fmt.Printf("Failed to start NMI profiling\n")
		}
	} else {
		// stop profiling
		rips, isbt := profhw.stopnmi()
		if len(rips) == 0 {
			fmt.Printf("No samples!\n")
			return
		}
		fmt.Printf("%v samples\n", len(rips))

		if isbt {
			pd := &fs.Profdev
			pd.Lock()
			pd.Prips.Reset()
			pd.Bts = rips
			pd.Unlock()
		} else {
			m := make(map[uintptr]int)
			for _, v := range rips {
				m[v] = m[v] + 1
			}
			prips := fs.Perfrips_t{}
			prips.Init(m)
			sort.Sort(sort.Reverse(&prips))

			pd := &fs.Profdev
			pd.Lock()
			pd.Prips = prips
			pd.Bts = nil
			pd.Unlock()
		}
	}
}

var hacklock sync.Mutex
var hackctrs []int

func _prof_pmc(en bool, events []pmev_t) {
	hacklock.Lock()
	defer hacklock.Unlock()

	if en {
		if hackctrs != nil {
			fmt.Printf("counters in use\n")
			return
		}
		cs, ok := profhw.startpmc(events)
		if ok {
			hackctrs = cs
		} else {
			fmt.Printf("failed to start counters\n")
		}
	} else {
		if hackctrs == nil {
			return
		}
		r := profhw.stoppmc(hackctrs)
		hackctrs = nil
		for i, ev := range events {
			t := ""
			if ev.pflags&EVF_USR != 0 {
				t = "(usr"
			}
			if ev.pflags&EVF_OS != 0 {
				if t != "" {
					t += "+os"
				} else {
					t = "(os"
				}
			}
			if t != "" {
				t += ")"
			}
			n := pmevid_names[ev.evid] + " " + t
			fmt.Printf("%-30s: %15v\n", n, r[i])
		}
	}
}

var fakeptr *common.Proc_t

//var fakedur = make([][]uint8, 256)
//var duri int

func sys_prof(proc *common.Proc_t, ptype, _events, _pmflags, intperiod int) int {
	en := true
	if ptype&common.PROF_DISABLE != 0 {
		en = false
	}
	pmflags := pmflag_t(_pmflags)
	switch {
	case ptype&common.PROF_GOLANG != 0:
		_prof_go(en)
	case ptype&common.PROF_SAMPLE != 0:
		ev := pmev_t{evid: pmevid_t(_events),
			pflags: pmflags}
		_prof_nmi(en, ev, intperiod)
	case ptype&common.PROF_COUNT != 0:
		if pmflags&EVF_BACKTRACE != 0 {
			return int(-defs.EINVAL)
		}
		evs := make([]pmev_t, 0, 4)
		for i := uint(0); i < 64; i++ {
			b := 1 << i
			if _events&b != 0 {
				n := pmev_t{}
				n.evid = pmevid_t(b)
				n.pflags = pmflags
				evs = append(evs, n)
			}
		}
		_prof_pmc(en, evs)
	case ptype&common.PROF_HACK != 0:
		runtime.Setheap(_events << 20)
	case ptype&common.PROF_HACK2 != 0:
		if _events < 0 {
			return int(-defs.EINVAL)
		}
		fmt.Printf("GOGC = %v\n", _events)
		debug.SetGCPercent(_events)
	case ptype&common.PROF_HACK3 != 0:
		if _events < 0 {
			return int(-defs.EINVAL)
		}
		buf := make([]uint8, _events)
		if buf == nil {
		}
		//fakedur[duri] = buf
		//duri = (duri + 1) % len(fakedur)
		//for i := 0; i < _events/8; i++ {
		//fakeptr = proc
		//}
	case ptype&common.PROF_HACK4 != 0:
		makefake(proc)
	case ptype&common.PROF_HACK5 != 0:
		n := _events
		if n < 0 {
			return int(-defs.EINVAL)
		}
		runtime.SetMaxheap(n)
		fmt.Printf("remaining mem: %v\n",
			res.Human(runtime.Memremain()))
	default:
		return int(-defs.EINVAL)
	}
	return 0
}

func makefake(p *common.Proc_t) defs.Err_t {
	p.Fdl.Lock()
	defer p.Fdl.Unlock()

	made := 0
	const want = 1e6
	newfds := make([]*vm.Fd_t, want)

	for times := 0; times < 4; times++ {
		fmt.Printf("close half...\n")
		for i := 0; i < len(newfds)/2; i++ {
			newfds[i] = nil
		}
		// sattolos
		for i := len(newfds) - 1; i >= 0; i-- {
			si := rand.Intn(i + 1)
			t := newfds[i]
			newfds[i] = newfds[si]
			newfds[si] = t
		}
		for i := range newfds {
			if newfds[i] == nil {
				newfds[i] = thefs.Makefake()
			}
		}
	}

	for i := range newfds {
		if i < len(p.Fds) && p.Fds[i] != nil {
			newfds[i] = p.Fds[i]
		} else {
			made++
		}
	}
	p.Fds = newfds
	fmt.Printf("bloat finished %v\n", made)
	return 0
}

func sys_info(proc *common.Proc_t, n int) int {
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)

	ret := int(-defs.EINVAL)
	switch n {
	case common.SINFO_GCCOUNT:
		ret = int(ms.NumGC)
	case common.SINFO_GCPAUSENS:
		ret = int(ms.PauseTotalNs)
	case common.SINFO_GCHEAPSZ:
		ret = int(ms.Alloc)
		fmt.Printf("Total heap size: %v MB (%v MB)\n",
			runtime.Heapsz()/(1<<20), ms.Alloc>>20)
	case common.SINFO_GCMS:
		tot := runtime.GCmarktime() + runtime.GCbgsweeptime()
		ret = tot / 1000000
	case common.SINFO_GCTOTALLOC:
		ret = int(ms.TotalAlloc)
	case common.SINFO_GCMARKT:
		ret = runtime.GCmarktime() / 1000000
	case common.SINFO_GCSWEEPT:
		ret = runtime.GCbgsweeptime() / 1000000
	case common.SINFO_GCWBARRT:
		ret = runtime.GCwbenabledtime() / 1000000
	case common.SINFO_GCOBJS:
		ret = int(ms.HeapObjects)
	case 10:
		runtime.GC()
		ret = 0
		p1, p2 := physmem.Pgcount()
		fmt.Printf("pgcount: %v, %v\n", p1, p2)
	case 11:
		//proc.Aspace.Vmregion.dump()
		fmt.Printf("proc dump:\n")
		common.Proclock.Lock()
		for i := range common.Allprocs {
			fmt.Printf("   %3v %v\n", common.Allprocs[i].Pid, common.Allprocs[i].Name)
		}
		common.Proclock.Unlock()
		ret = 0
	}

	return ret
}

func readn(a []uint8, n int, off int) int {
	p := unsafe.Pointer(&a[off])
	var ret int
	switch n {
	case 8:
		ret = *(*int)(p)
	case 4:
		ret = int(*(*uint32)(p))
	case 2:
		ret = int(*(*uint16)(p))
	case 1:
		ret = int(*(*uint8)(p))
	default:
		panic("no")
	}
	return ret
}

func writen(a []uint8, sz int, off int, val int) {
	p := unsafe.Pointer(&a[off])
	switch sz {
	case 8:
		*(*int)(p) = val
	case 4:
		*(*uint32)(p) = uint32(val)
	case 2:
		*(*uint16)(p) = uint16(val)
	case 1:
		*(*uint8)(p) = uint8(val)
	default:
		panic("no")
	}
}

// returns the byte size/offset of field n. they can be used to read []data.
func fieldinfo(sizes []int, n int) (int, int) {
	if n >= len(sizes) {
		panic("bad field number")
	}
	off := 0
	for i := 0; i < n; i++ {
		off += sizes[i]
	}
	return sizes[n], off
}

type elf_t struct {
	data []uint8
}

type elf_phdr struct {
	etype   int
	flags   int
	vaddr   int
	filesz  int
	fileoff int
	memsz   int
}

const (
	ELF_QUARTER = 2
	ELF_HALF    = 4
	ELF_OFF     = 8
	ELF_ADDR    = 8
	ELF_XWORD   = 8
)

func (e *elf_t) sanity() bool {
	// make sure its an elf
	e_ident := 0
	elfmag := 0x464c457f
	t := readn(e.data, ELF_HALF, e_ident)
	if t != elfmag {
		return false
	}

	// and that we read the entire elf header and program headers
	dlen := len(e.data)

	e_ehsize := 0x34
	ehlen := readn(e.data, ELF_QUARTER, e_ehsize)
	if dlen < ehlen {
		fmt.Printf("read too few elf bytes (elf header)\n")
		return false
	}

	e_phoff := 0x20
	e_phentsize := 0x36
	e_phnum := 0x38

	poff := readn(e.data, ELF_OFF, e_phoff)
	phsz := readn(e.data, ELF_QUARTER, e_phentsize)
	phnum := readn(e.data, ELF_QUARTER, e_phnum)
	phend := poff + phsz*phnum
	if dlen < phend {
		fmt.Printf("read too few elf bytes (program headers)\n")
		return false
	}

	return true
}

func (e *elf_t) npheaders() int {
	e_phnum := 0x38
	return readn(e.data, ELF_QUARTER, e_phnum)
}

func (e *elf_t) header(c int) elf_phdr {
	ret := elf_phdr{}

	nph := e.npheaders()
	if c >= nph {
		panic("header idx too large")
	}
	d := e.data
	e_phoff := 0x20
	e_phentsize := 0x36
	hoff := readn(d, ELF_OFF, e_phoff)
	hsz := readn(d, ELF_QUARTER, e_phentsize)

	p_type := 0x0
	p_flags := 0x4
	p_offset := 0x8
	p_vaddr := 0x10
	p_filesz := 0x20
	p_memsz := 0x28
	f := func(w int, sz int) int {
		return readn(d, sz, hoff+c*hsz+w)
	}
	ret.etype = f(p_type, ELF_HALF)
	ret.flags = f(p_flags, ELF_HALF)
	ret.fileoff = f(p_offset, ELF_OFF)
	ret.vaddr = f(p_vaddr, ELF_ADDR)
	ret.filesz = f(p_filesz, ELF_XWORD)
	ret.memsz = f(p_memsz, ELF_XWORD)
	return ret
}

func (e *elf_t) headers() []elf_phdr {
	pnum := e.npheaders()
	ret := make([]elf_phdr, pnum)
	for i := 0; i < pnum; i++ {
		ret[i] = e.header(i)
	}
	return ret
}

func (e *elf_t) entry() int {
	e_entry := 0x18
	return readn(e.data, ELF_ADDR, e_entry)
}

func segload(proc *common.Proc_t, entry int, hdr *elf_phdr, fops vm.Fdops_i) defs.Err_t {
	if hdr.vaddr%mem.PGSIZE != hdr.fileoff%mem.PGSIZE {
		panic("requires copying")
	}
	perms := vm.PTE_U
	//PF_X := 1
	PF_W := 2
	if hdr.flags&PF_W != 0 {
		perms |= vm.PTE_W
	}

	var did int
	// the bss segment's virtual address may start on the same page as the
	// previous segment. if that is the case, we may not be able to avoid
	// copying.
	// XXX this doesn't seem to happen anymore; why was it ever the case?
	if _, ok := proc.Aspace.Vmregion.Lookup(uintptr(hdr.vaddr)); ok {
		panic("htf")
		va := hdr.vaddr
		pg, err := proc.Aspace.Userdmap8_inner(va, true)
		if err != 0 {
			return err
		}
		mmapi, err := fops.Mmapi(hdr.fileoff, 1, false)
		if err != 0 {
			return err
		}
		bsrc := mem.Pg2bytes(mmapi[0].Pg)[:]
		bsrc = bsrc[va&int(vm.PGOFFSET):]
		if len(pg) > hdr.filesz {
			pg = pg[0:hdr.filesz]
		}
		copy(pg, bsrc)
		did = len(pg)
	}
	filesz := util.Roundup(hdr.vaddr+hdr.filesz-did, mem.PGSIZE)
	filesz -= util.Rounddown(hdr.vaddr, mem.PGSIZE)
	proc.Aspace.Vmadd_file(hdr.vaddr+did, filesz, perms, fops, hdr.fileoff+did)
	// eagerly map the page at the entry address
	if entry >= hdr.vaddr && entry < hdr.vaddr+hdr.memsz {
		ent := uintptr(entry)
		vmi, ok := proc.Aspace.Vmregion.Lookup(ent)
		if !ok {
			panic("just mapped?")
		}
		err := vm.Sys_pgfault(&proc.Aspace, vmi, ent, uintptr(vm.PTE_U))
		if err != 0 {
			return err
		}
	}
	if hdr.filesz == hdr.memsz {
		return 0
	}
	// the bss must be zero, but the first bss address may lie on a page
	// which is mapped into the page cache. thus we must create a
	// per-process copy and zero the bss bytes in the copy.
	bssva := hdr.vaddr + hdr.filesz
	bsslen := hdr.memsz - hdr.filesz
	if bssva&int(vm.PGOFFSET) != 0 {
		bpg, err := proc.Aspace.Userdmap8_inner(bssva, true)
		if err != 0 {
			return err
		}
		if bsslen < len(bpg) {
			bpg = bpg[:bsslen]
		}
		copy(bpg, mem.Zerobpg[:])
		bssva += len(bpg)
		bsslen = util.Roundup(bsslen-len(bpg), mem.PGSIZE)
	}
	// bss may have been completely contained in the copied page.
	if bsslen > 0 {
		proc.Aspace.Vmadd_anon(bssva, util.Roundup(bsslen, mem.PGSIZE), perms)
	}
	return 0
}

// returns user address of read-only TLS, thread 0's TLS image, TLS size, and
// success. caller must hold proc's pagemap lock.
func (e *elf_t) elf_load(proc *common.Proc_t, f *vm.Fd_t) (int, int, int, defs.Err_t) {
	PT_LOAD := 1
	PT_TLS := 7
	istls := false
	tlssize := 0
	var tlsaddr int
	var tlscopylen int

	gimme := bounds.Bounds(bounds.B_ELF_T_ELF_LOAD)
	entry := e.entry()
	// load each elf segment directly into process memory
	for _, hdr := range e.headers() {
		// XXX get rid of worthless user program segments
		if !res.Resadd_noblock(gimme) {
			return 0, 0, 0, -defs.ENOHEAP
		}
		if hdr.etype == PT_TLS {
			istls = true
			tlsaddr = hdr.vaddr
			tlssize = util.Roundup(hdr.memsz, 8)
			tlscopylen = hdr.filesz
		} else if hdr.etype == PT_LOAD && hdr.vaddr >= mem.USERMIN {
			err := segload(proc, entry, &hdr, f.Fops)
			if err != 0 {
				return 0, 0, 0, err
			}
		}
	}

	freshtls := 0
	t0tls := 0
	if istls {
		// create fresh TLS image and map it COW for thread 0
		l := util.Roundup(tlsaddr+tlssize, mem.PGSIZE)
		l -= util.Rounddown(tlsaddr, mem.PGSIZE)

		freshtls = proc.Aspace.Unusedva_inner(0, 2*l)
		t0tls = freshtls + l
		proc.Aspace.Vmadd_anon(freshtls, l, vm.PTE_U)
		proc.Aspace.Vmadd_anon(t0tls, l, vm.PTE_U|vm.PTE_W)
		perms := vm.PTE_U

		for i := 0; i < l; i += mem.PGSIZE {
			// allocator zeros objects, so tbss is already
			// initialized.
			_, p_pg, ok := physmem.Refpg_new()
			if !ok {
				return 0, 0, 0, -defs.ENOMEM
			}
			_, ok = proc.Aspace.Page_insert(freshtls+i, p_pg, perms,
				true)
			if !ok {
				physmem.Refdown(p_pg)
				return 0, 0, 0, -defs.ENOMEM
			}
			// map fresh TLS for thread 0
			nperms := perms | vm.PTE_COW
			_, ok = proc.Aspace.Page_insert(t0tls+i, p_pg, nperms, true)
			if !ok {
				physmem.Refdown(p_pg)
				return 0, 0, 0, -defs.ENOMEM
			}
		}
		// copy TLS data to freshtls
		tlsvmi, ok := proc.Aspace.Vmregion.Lookup(uintptr(tlsaddr))
		if !ok {
			panic("must succeed")
		}
		for i := 0; i < tlscopylen; {
			if !res.Resadd_noblock(gimme) {
				return 0, 0, 0, -defs.ENOHEAP
			}

			_src, p_pg, err := tlsvmi.Filepage(uintptr(tlsaddr + i))
			if err != 0 {
				return 0, 0, 0, err
			}
			off := (tlsaddr + i) & int(vm.PGOFFSET)
			src := mem.Pg2bytes(_src)[off:]
			bpg, err := proc.Aspace.Userdmap8_inner(freshtls+i, true)
			if err != 0 {
				physmem.Refdown(p_pg)
				return 0, 0, 0, err
			}
			left := tlscopylen - i
			if len(src) > left {
				src = src[0:left]
			}
			copy(bpg, src)
			i += len(src)
			physmem.Refdown(p_pg)
		}

		// amd64 sys 5 abi specifies that the tls pointer references to
		// the first invalid word past the end of the tls
		t0tls += tlssize
	}
	return freshtls, t0tls, tlssize, 0
}
