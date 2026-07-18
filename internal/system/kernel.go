// internal/system/kernel.go
package system

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var kernelFlavorSuffixRegexp = regexp.MustCompile(`-[a-z]+$`)

// parseKernelBase strips a trailing "-<flavor>" suffix from a `uname -r`
// style string, e.g. "5.15.0-91-generic" -> "5.15.0-91".
func parseKernelBase(unameR string) string {
	return kernelFlavorSuffixRegexp.ReplaceAllString(unameR, "")
}

// parseInstalledKernels parses `dpkg --list 'linux-image-*'` output,
// returning package names for lines actually installed ("ii" status).
func parseInstalledKernels(dpkgOutput string) []string {
	var pkgs []string
	for _, line := range strings.Split(dpkgOutput, "\n") {
		if !strings.HasPrefix(line, "ii") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pkgs = append(pkgs, fields[1])
	}
	return pkgs
}

var kernelImageRegexp = regexp.MustCompile(`^linux-image-[0-9]+\.[0-9]+`)

// oldKernels filters installedPkgs down to versioned kernel image packages
// (excluding meta-packages like "linux-image-generic") that don't belong
// to the currently running kernel, matched by kernelBase substring -
// mirroring the bash reference's filtering exactly.
func oldKernels(installedPkgs []string, kernelBase string) []string {
	var old []string
	for _, pkg := range installedPkgs {
		if !kernelImageRegexp.MatchString(pkg) {
			continue
		}
		if strings.Contains(pkg, kernelBase) {
			continue
		}
		old = append(old, pkg)
	}
	return old
}

// classifyKernelRemoval decides Warning vs Unsafe for a kernel-removal
// action: Unsafe only if the exact currently-BOOTED kernel string (not
// just its base) still appears in the candidate list - a defensive
// second check beyond oldKernels' filtering, reproduced from the bash
// reference exactly (including its known limitation: this is defense in
// depth, not a hard guarantee).
func classifyKernelRemoval(old []string, currentKernel string) Tier {
	for _, pkg := range old {
		if strings.Contains(pkg, currentKernel) {
			return TierUnsafe
		}
	}
	return TierWarning
}

// analyzeKernels inspects the running kernel and installed kernel
// packages and returns a kernel-removal Action if there's anything to
// remove, or nil if there's nothing to do. Purely read-only (`uname -r`,
// `dpkg --list`) - never mutates anything.
func analyzeKernels() (*Action, error) {
	unameOut, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return nil, fmt.Errorf("checking current kernel: %w", err)
	}
	currentKernel := strings.TrimSpace(string(unameOut))
	kernelBase := parseKernelBase(currentKernel)

	dpkgOut, err := exec.Command("dpkg", "--list", "linux-image-*").Output()
	if err != nil {
		// dpkg --list can exit non-zero when the pattern matches nothing -
		// treat as "no kernel packages found", not a fatal error.
		return nil, nil
	}
	installed := parseInstalledKernels(string(dpkgOut))
	old := oldKernels(installed, kernelBase)
	if len(old) == 0 {
		return nil, nil
	}

	tier := classifyKernelRemoval(old, currentKernel)
	reason := fmt.Sprintf("Kernels to remove: %s", strings.Join(old, ", "))
	if tier == TierUnsafe {
		reason = "CRITICAL: current kernel detected in removal list!"
	}
	return &Action{
		Name:    "Remove old kernels",
		Tier:    tier,
		Command: append([]string{"apt-get", "purge", "-y"}, old...),
		Reason:  reason,
	}, nil
}
