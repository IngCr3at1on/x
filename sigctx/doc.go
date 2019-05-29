// Package sigctx provides cancel signal wrapped context.
// This is a super light package which I find myself using in most of my projects
// rather than continuing to copy/paste this one segment of code across a million
// repos I decided to just stick it here.
// A context can be cancelled with any of the following signals or os calls.
//
//	syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt and os.Kill
package sigctx
