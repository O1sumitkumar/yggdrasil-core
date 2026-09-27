package tools

import "testing"

func TestPolicyDeny(t *testing.T) {
	p := NewPolicyEngine()
	allowed, prompt, err := p.Decide("terminal", PolicyDeny)
	if err == nil || allowed || prompt {
		t.Fatalf("deny: allowed=%v prompt=%v err=%v", allowed, prompt, err)
	}
}

func TestPolicyAllow(t *testing.T) {
	p := NewPolicyEngine()
	allowed, prompt, err := p.Decide("git.status", PolicyAllow)
	if err != nil || !allowed || prompt {
		t.Fatalf("allow: allowed=%v prompt=%v err=%v", allowed, prompt, err)
	}
}

func TestPolicyAllowForSession(t *testing.T) {
	p := NewPolicyEngine()
	allowed, prompt, _ := p.Decide("filesystem.write", PolicyAllowForSession)
	if allowed || !prompt {
		t.Fatal("expected prompt first")
	}
	p.AllowSession("filesystem.write")
	allowed, prompt, err := p.Decide("filesystem.write", PolicyAllowForSession)
	if err != nil || !allowed || prompt {
		t.Fatalf("session: allowed=%v prompt=%v err=%v", allowed, prompt, err)
	}
}

func TestPolicyAsk(t *testing.T) {
	p := NewPolicyEngine()
	allowed, prompt, err := p.Decide("terminal", PolicyAsk)
	if err != nil || allowed || !prompt {
		t.Fatalf("ask: allowed=%v prompt=%v err=%v", allowed, prompt, err)
	}
}
