package ethnonce

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/garyburd/redigo/redis"
	"testing"
)

func _testinit() (redis.Conn, *ethclient.Client) {
	conn, _ := ethclient.Dial("http://localhost:18545")
	c, _ := redis.Dial("tcp", ":6379")
	return c, conn
}

func TestGiveCommit(t *testing.T) {
	redis_conn, eth_conn := _testinit()
	mgr := NewNonceManager("testhash")
	addr := common.HexToAddress(`0xe35f3e2a93322b61e5d8931f806ff38f4a4f4d88`)
	mgr.SyncNonce(redis_conn, addr, eth_conn)
	nonce, err := mgr.MustGiveNonce(redis_conn, addr, eth_conn)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nonce:", nonce)
	if err = mgr.CommitNonce(redis_conn, addr, nonce, true); err != nil {
		t.Fatal(err)
	}
	nonce1, err := mgr.MustGiveNonce(redis_conn, addr)
	if err != nil {
		t.Fatal(err)
	}
	if nonce1 != nonce+1 {
		t.Fatal("bad ", nonce1)
	}
	_, err = mgr.GiveNonce(redis_conn, addr)
	if err == nil {
		t.Fatal("should err")
	}
	if err = mgr.CommitNonce(redis_conn, addr, nonce1, false); err != nil {
		t.Fatal(err)
	}
}
func TestSync(t *testing.T) {
	redis_conn, eth_conn := _testinit()
	mgr := NewNonceManager("testhash1")
	addr := common.HexToAddress(`0xe35f3e2a93322b61e5d8931f806ff38f4a4f4d88`)
	if _, err := mgr.SyncNonce(redis_conn, addr, eth_conn); err != nil {
		t.Fatal("first sync should OK,but ", err)
	}
	nonce, err := mgr.MustGiveNonce(redis_conn, addr, eth_conn)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("get nonce ", nonce)
	if _, err := mgr.SyncNonce(redis_conn, addr, eth_conn); err == nil {
		t.Fatal("second sync should fail")
	}
	t.Log("should sync failed,pass")
}
