package redsrv

import (
	"context"
	"github.com/go-redis/redis/v7"
	"time"
)

func (m *MockCmdable) Pipeline() redis.Pipeliner {
	panic("implement me")
}

func (m *MockCmdable) Pipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	panic("implement me")
}

func (m *MockCmdable) TxPipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	panic("implement me")
}

func (m *MockCmdable) TxPipeline() redis.Pipeliner {
	panic("implement me")
}

func (m *MockCmdable) Command() *redis.CommandsInfoCmd {
	panic("implement me")
}

func (m *MockCmdable) ClientGetName() *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) Echo(message interface{}) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) Ping() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) Quit() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) Del(keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) Unlink(keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) Dump(key string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) Exists(keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) Keys(pattern string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) Migrate(host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) Move(key string, db int) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) ObjectRefCount(key string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ObjectEncoding(key string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) ObjectIdleTime(key string) *redis.DurationCmd {
	panic("implement me")
}

func (m *MockCmdable) Persist(key string) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) PExpire(key string, expiration time.Duration) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) PExpireAt(key string, tm time.Time) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) PTTL(key string) *redis.DurationCmd {
	panic("implement me")
}

func (m *MockCmdable) RandomKey() *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) Rename(key, newkey string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) RenameNX(key, newkey string) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) Restore(key string, ttl time.Duration, value string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) RestoreReplace(key string, ttl time.Duration, value string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) Sort(key string, sort *redis.Sort) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) SortStore(key, store string, sort *redis.Sort) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) SortInterfaces(key string, sort *redis.Sort) *redis.SliceCmd {
	panic("implement me")
}

func (m *MockCmdable) Touch(keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) TTL(key string) *redis.DurationCmd {
	panic("implement me")
}

func (m *MockCmdable) Type(key string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	panic("implement me")
}

func (m *MockCmdable) SScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	panic("implement me")
}

func (m *MockCmdable) HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	panic("implement me")
}

func (m *MockCmdable) ZScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	panic("implement me")
}

func (m *MockCmdable) Append(key, value string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) BitCount(key string, bitCount *redis.BitCount) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) BitOpAnd(destKey string, keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) BitOpOr(destKey string, keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) BitOpXor(destKey string, keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) BitOpNot(destKey string, key string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) BitPos(key string, bit int64, pos ...int64) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) BitField(key string, args ...interface{}) *redis.IntSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) Decr(key string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) DecrBy(key string, decrement int64) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) Get(key string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) GetBit(key string, offset int64) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) GetRange(key string, start, end int64) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) GetSet(key string, value interface{}) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) Incr(key string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) IncrBy(key string, value int64) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) IncrByFloat(key string, value float64) *redis.FloatCmd {
	panic("implement me")
}

func (m *MockCmdable) MGet(keys ...string) *redis.SliceCmd {
	panic("implement me")
}

func (m *MockCmdable) MSet(values ...interface{}) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) MSetNX(values ...interface{}) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) SetBit(key string, offset int64, value int) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) SetRange(key string, offset int64, value string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) StrLen(key string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) HDel(key string, fields ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) HExists(key, field string) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) HGet(key, field string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) HGetAll(key string) *redis.StringStringMapCmd {
	panic("implement me")
}

func (m *MockCmdable) HIncrBy(key, field string, incr int64) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) HIncrByFloat(key, field string, incr float64) *redis.FloatCmd {
	panic("implement me")
}

func (m *MockCmdable) HKeys(key string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) HLen(key string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) HMGet(key string, fields ...string) *redis.SliceCmd {
	panic("implement me")
}

func (m *MockCmdable) HMSet(key string, values ...interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) HSet(key, field string, value interface{}) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) HSetNX(key, field string, value interface{}) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) HVals(key string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) BRPopLPush(source, destination string, timeout time.Duration) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) LIndex(key string, index int64) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) LInsert(key, op string, pivot, value interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) LInsertBefore(key string, pivot, value interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) LInsertAfter(key string, pivot, value interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) LLen(key string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) LPop(key string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) LPush(key string, values ...interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) LPushX(key string, values ...interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) LRem(key string, count int64, value interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) LSet(key string, index int64, value interface{}) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) LTrim(key string, start, stop int64) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) RPop(key string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) RPopLPush(source, destination string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) RPush(key string, values ...interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) RPushX(key string, values ...interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) SAdd(key string, members ...interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) SCard(key string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) SDiff(keys ...string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) SDiffStore(destination string, keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) SInter(keys ...string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) SInterStore(destination string, keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) SIsMember(key string, member interface{}) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) SMembers(key string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) SMembersMap(key string) *redis.StringStructMapCmd {
	panic("implement me")
}

func (m *MockCmdable) SMove(source, destination string, member interface{}) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) SPop(key string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) SPopN(key string, count int64) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) SRandMember(key string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) SRandMemberN(key string, count int64) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) SRem(key string, members ...interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) SUnion(keys ...string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) SUnionStore(destination string, keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) XAdd(a *redis.XAddArgs) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) XDel(stream string, ids ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) XLen(stream string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) XRange(stream, start, stop string) *redis.XMessageSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) XRangeN(stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) XRevRange(stream string, start, stop string) *redis.XMessageSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) XRevRangeN(stream string, start, stop string, count int64) *redis.XMessageSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) XRead(a *redis.XReadArgs) *redis.XStreamSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) XReadStreams(streams ...string) *redis.XStreamSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) XGroupCreate(stream, group, start string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) XGroupCreateMkStream(stream, group, start string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) XGroupSetID(stream, group, start string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) XGroupDestroy(stream, group string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) XGroupDelConsumer(stream, group, consumer string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) XReadGroup(a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) XAck(stream, group string, ids ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) XPending(stream, group string) *redis.XPendingCmd {
	panic("implement me")
}

func (m *MockCmdable) XPendingExt(a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	panic("implement me")
}

func (m *MockCmdable) XClaim(a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) XClaimJustID(a *redis.XClaimArgs) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) XTrim(key string, maxLen int64) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) XTrimApprox(key string, maxLen int64) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) XInfoGroups(key string) *redis.XInfoGroupsCmd {
	panic("implement me")
}

func (m *MockCmdable) BZPopMax(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	panic("implement me")
}

func (m *MockCmdable) BZPopMin(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	panic("implement me")
}

func (m *MockCmdable) ZAdd(key string, members ...*redis.Z) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZAddNX(key string, members ...*redis.Z) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZAddXX(key string, members ...*redis.Z) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZAddCh(key string, members ...*redis.Z) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZAddNXCh(key string, members ...*redis.Z) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZAddXXCh(key string, members ...*redis.Z) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZIncr(key string, member *redis.Z) *redis.FloatCmd {
	panic("implement me")
}

func (m *MockCmdable) ZIncrNX(key string, member *redis.Z) *redis.FloatCmd {
	panic("implement me")
}

func (m *MockCmdable) ZIncrXX(key string, member *redis.Z) *redis.FloatCmd {
	panic("implement me")
}

func (m *MockCmdable) ZCard(key string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZCount(key, min, max string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZLexCount(key, min, max string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZIncrBy(key string, increment float64, member string) *redis.FloatCmd {
	panic("implement me")
}

func (m *MockCmdable) ZInterStore(destination string, store *redis.ZStore) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZPopMax(key string, count ...int64) *redis.ZSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZPopMin(key string, count ...int64) *redis.ZSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRange(key string, start, stop int64) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRangeByScore(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRangeByLex(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRangeByScoreWithScores(key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRank(key, member string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRem(key string, members ...interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRemRangeByRank(key string, start, stop int64) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRemRangeByScore(key, min, max string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRemRangeByLex(key, min, max string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRevRange(key string, start, stop int64) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRevRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRevRangeByScore(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRevRangeByLex(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ZRevRank(key, member string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ZScore(key, member string) *redis.FloatCmd {
	panic("implement me")
}

func (m *MockCmdable) ZUnionStore(dest string, store *redis.ZStore) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) PFAdd(key string, els ...interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) PFCount(keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) PFMerge(dest string, keys ...string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) BgRewriteAOF() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) BgSave() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClientKill(ipPort string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClientKillByFilter(keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ClientList() *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) ClientPause(dur time.Duration) *redis.BoolCmd {
	panic("implement me")
}

func (m *MockCmdable) ClientID() *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ConfigGet(parameter string) *redis.SliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ConfigResetStat() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ConfigSet(parameter, value string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ConfigRewrite() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) DBSize() *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) FlushAll() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) FlushAllAsync() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) FlushDB() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) FlushDBAsync() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) Info(section ...string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) LastSave() *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) Save() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) Shutdown() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ShutdownSave() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ShutdownNoSave() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) SlaveOf(host, port string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) Time() *redis.TimeCmd {
	panic("implement me")
}

func (m *MockCmdable) Eval(script string, keys []string, args ...interface{}) *redis.Cmd {
	panic("implement me")
}

func (m *MockCmdable) EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	panic("implement me")
}

func (m *MockCmdable) ScriptExists(hashes ...string) *redis.BoolSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ScriptFlush() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ScriptKill() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ScriptLoad(script string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) DebugObject(key string) *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) Publish(channel string, message interface{}) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) PubSubChannels(pattern string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) PubSubNumSub(channels ...string) *redis.StringIntMapCmd {
	panic("implement me")
}

func (m *MockCmdable) PubSubNumPat() *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterSlots() *redis.ClusterSlotsCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterNodes() *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterMeet(host, port string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterForget(nodeID string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterReplicate(nodeID string) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterResetSoft() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterResetHard() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterInfo() *redis.StringCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterKeySlot(key string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterGetKeysInSlot(slot int, count int) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterCountFailureReports(nodeID string) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterCountKeysInSlot(slot int) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterDelSlots(slots ...int) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterDelSlotsRange(min, max int) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterSaveConfig() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterSlaves(nodeID string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterFailover() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterAddSlots(slots ...int) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ClusterAddSlotsRange(min, max int) *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) GeoAdd(key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) GeoPos(key string, members ...string) *redis.GeoPosCmd {
	panic("implement me")
}

func (m *MockCmdable) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	panic("implement me")
}

func (m *MockCmdable) GeoRadiusStore(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) GeoRadiusByMember(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	panic("implement me")
}

func (m *MockCmdable) GeoRadiusByMemberStore(key, member string, query *redis.GeoRadiusQuery) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) GeoDist(key string, member1, member2, unit string) *redis.FloatCmd {
	panic("implement me")
}

func (m *MockCmdable) GeoHash(key string, members ...string) *redis.StringSliceCmd {
	panic("implement me")
}

func (m *MockCmdable) ReadOnly() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) ReadWrite() *redis.StatusCmd {
	panic("implement me")
}

func (m *MockCmdable) MemoryUsage(key string, samples ...int) *redis.IntCmd {
	panic("implement me")
}

func (m *MockCmdable) Context() context.Context {
	panic("implement me")
}

func (m *MockCmdable) AddHook(redis.Hook) {
	panic("implement me")
}

func (m *MockCmdable) Watch(fn func(*redis.Tx) error, keys ...string) error {
	panic("implement me")
}

func (m *MockCmdable) Do(args ...interface{}) *redis.Cmd {
	panic("implement me")
}

func (m *MockCmdable) DoContext(ctx context.Context, args ...interface{}) *redis.Cmd {
	panic("implement me")
}

func (m *MockCmdable) Process(cmd redis.Cmder) error {
	panic("implement me")
}

func (m *MockCmdable) ProcessContext(ctx context.Context, cmd redis.Cmder) error {
	panic("implement me")
}

func (m *MockCmdable) Subscribe(channels ...string) *redis.PubSub {
	panic("implement me")
}

func (m *MockCmdable) PSubscribe(channels ...string) *redis.PubSub {
	panic("implement me")
}

func (m *MockCmdable) Close() error {
	panic("implement me")
}
