package main

import (
	"testing"
	"time"

	"github.com/buildkite/interpolate"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/ptr"

	"github.com/delamonpansie/silencer/config"
	"github.com/delamonpansie/silencer/filter"
)

func Test_banWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	blocker := filter.NewMockBlocker(ctrl)

	blocker.EXPECT().List().Return(nil)

	blocker.EXPECT().Block("ip1")
	blocker.EXPECT().Block("ip2")
	blocker.EXPECT().Block("ip3")

	blocker.EXPECT().Unblock("ip1")
	blocker.EXPECT().Unblock("ip3")
	blocker.EXPECT().Unblock("ip2")

	c := worker(blocker, time.Second)
	c <- blockRequest{"ip1", time.Second}
	c <- blockRequest{"ip1", time.Second}
	c <- blockRequest{"ip2", 3 * time.Second}
	c <- blockRequest{"ip3", time.Second}

	time.Sleep(5 * time.Second)
}

func testRule(re ...string) rule {
	cfg := config.Load()
	env := interpolate.NewMapEnv(cfg.Env)
	return newRule("testRule", re, time.Second, env)
}

func Test_rule1(t *testing.T) {
	rule := testRule(
		`^$date_time \[\d+\] SMTP protocol error in "AUTH LOGIN" (.*)`,
		`(.*) AUTH command used when not advertised$`,
		`^H=\(\S+\) \[($ip)\]`,
	)

	line := `2020-10-27 21:05:50.780 [2168] SMTP protocol error in "AUTH LOGIN" H=(User) [103.253.42.54]:57715 I=[85.218.130.46]:25 AUTH command used when not advertised
`
	m, err := rule.match(line)
	require.NoError(t, err)
	assert.Equal(t, ptr.String("103.253.42.54"), m)
}

func Test_rule2(t *testing.T) {
	rule := testRule("aaa", `($ip)`)
	m, err := rule.match("aaa bbb 1.1.1.1")
	require.NoError(t, err)
	assert.Equal(t, ptr.String("1.1.1.1"), m)
}

func Test_rule22(t *testing.T) {
	rule := testRule("aaa (.*)", " (.*)", `$ip`)
	m, err := rule.match("aaa bbb 1.1.1.1")
	require.NoError(t, err)
	assert.Equal(t, ptr.String("1.1.1.1"), m)
}

func Test_rule3(t *testing.T) {
	rule := testRule("zzz", `($ip)`)
	m, err := rule.match("aaa bbb 1.1.1.1")
	require.NoError(t, err)
	assert.Nil(t, m)
}
