package storage

import (
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	testAK       = os.Getenv("QINIU_ACCESS_KEY")
	testSK       = os.Getenv("QINIU_SECRET_KEY")
	testBucket   = os.Getenv("QINIU_TEST_BUCKET")
	testKey      = "qiniu.png"
	testFetchUrl = "http://devtools.qiniu.com/qiniu.png"
	testSiteUrl  = "http://devtools.qiniu.com"
)

var mac *qbox.Mac
var bucketManager *BucketManager

func init() {
	if testAK == "" || testSK == "" {
		panic("please run ./test-env.sh first")
	}
	mac = qbox.NewMac(testAK, testSK)
	cfg := Config{}
	bucketManager = NewBucketManager(mac, &cfg)
	rand.Seed(time.Now().Unix())
}

//Test get zone
func TestGetZone(t *testing.T) {
	zone, err := GetZone(testAK, testBucket)
	if err != nil {
		t.Fatalf("GetZone() error, %s", err)
	}
	t.Log(zone.String())
}

//Test get bucket list
func TestBuckets(t *testing.T) {
	shared := true
	buckets, err := bucketManager.Buckets(shared)
	if err != nil {
		t.Fatalf("Buckets() error, %s", err)
	}

	for _, bucket := range buckets {
		t.Log(bucket)
	}
}

//Test get file info
func TestStat(t *testing.T) {
	keysToStat := []string{"qiniu.jpg"}

	for _, eachKey := range keysToStat {
		info, err := bucketManager.Stat(testBucket, eachKey)
		if err != nil {
			t.Logf("Stat() error, %s", err)
			t.Fail()
		} else {
			t.Logf("FileInfo:\n %s", info.String())
		}
	}
}

func TestCopyMoveDelete(t *testing.T) {
	keysCopyTarget := []string{"qiniu_1.jpg", "qiniu_2.jpg", "qiniu_3.jpg"}
	keysToDelete := make([]string, 0, len(keysCopyTarget))
	for _, eachKey := range keysCopyTarget {
		err := bucketManager.Copy(testBucket, testKey, testBucket, eachKey, true)
		if err != nil {
			t.Logf("Copy() error, %s", err)
			t.Fail()
		}
	}

	for _, eachKey := range keysCopyTarget {
		keyToDelete := eachKey + "_move"
		err := bucketManager.Move(testBucket, eachKey, testBucket, keyToDelete, true)
		if err != nil {
			t.Logf("Move() error, %s", err)
			t.Fail()
		} else {
			keysToDelete = append(keysToDelete, keyToDelete)
		}
	}

	for _, eachKey := range keysToDelete {
		err := bucketManager.Delete(testBucket, eachKey)
		if err != nil {
			t.Logf("Delete() error, %s", err)
			t.Fail()
		}
	}
}

func TestFetch(t *testing.T) {
	ret, err := bucketManager.Fetch(testFetchUrl, testBucket, "qiniu-fetch.png")
	if err != nil {
		t.Logf("Fetch() error, %s", err)
		t.Fail()
	} else {
		t.Logf("FetchRet:\n %s", ret.String())
	}
}

func TestFetchWithoutKey(t *testing.T) {
	ret, err := bucketManager.FetchWithoutKey(testFetchUrl, testBucket)
	if err != nil {
		t.Logf("FetchWithoutKey() error, %s", err)
		t.Fail()
	} else {
		t.Logf("FetchRet:\n %s", ret.String())
	}
}

func TestDeleteAfterDays(t *testing.T) {
	deleteKey := testKey + "_deleteAfterDays"
	days := 7
	bucketManager.Copy(testBucket, testKey, testBucket, deleteKey, true)
	err := bucketManager.DeleteAfterDays(testBucket, deleteKey, days)
	if err != nil {
		t.Logf("DeleteAfterDays() error, %s", err)
		t.Fail()
	}
}

func TestChangeMime(t *testing.T) {
	toChangeKey := testKey + "_changeMime"
	bucketManager.Copy(testBucket, testKey, testBucket, toChangeKey, true)
	newMime := "text/plain"
	err := bucketManager.ChangeMime(testBucket, toChangeKey, newMime)
	if err != nil {
		t.Fatalf("ChangeMime() error, %s", err)
	}

	info, err := bucketManager.Stat(testBucket, toChangeKey)
	if err != nil || info.MimeType != newMime {
		t.Fatalf("ChangeMime() failed, %s", err)
	}
	bucketManager.Delete(testBucket, toChangeKey)
}

func TestChangeType(t *testing.T) {
	toChangeKey := fmt.Sprintf("%s_changeType_%d", testKey, rand.Int())
	bucketManager.Copy(testBucket, testKey, testBucket, toChangeKey, true)
	fileType := 1
	err := bucketManager.ChangeType(testBucket, toChangeKey, fileType)
	if err != nil {
		t.Fatalf("ChangeType() error, %s", err)
	}

	info, err := bucketManager.Stat(testBucket, toChangeKey)
	if err != nil || info.Type != fileType {
		t.Fatalf("ChangeMime() failed, %s", err)
	}
	bucketManager.Delete(testBucket, toChangeKey)
}

func TestPrefetchAndImage(t *testing.T) {
	err := bucketManager.SetImage(testSiteUrl, testBucket)
	if err != nil {
		t.Fatalf("SetImage() error, %s", err)
	}

	err = bucketManager.Prefetch(testBucket, testKey)
	if err != nil {
		t.Fatalf("Prefetch() error, %s", err)
	}

	err = bucketManager.UnsetImage(testBucket)
	if err != nil {
		t.Fatalf("UnsetImage() error, %s", err)
	}
}

func TestListFiles(t *testing.T){
	limit:=100
	 prefix:="listfiles/"
	for i:=0;i<limit;i++{
		newKey:=fmt.Sprintf("%s%s/%d",prefix,testKey,i)
		bucketManager.Copy(testBucket,testKey,testBucket,newKey,true)
	}
	entries,_,_,hasNext,err:=bucketManager.ListFiles(testBucket,prefix,"","",limit)
	if err!=nil{
		t.Fatalf("ListFiles() error, %s",err)
	}

	if hasNext{
		t.Fatalf("ListFiles() failed, unexpected hasNext")
	}

	if len(entries)!=limit{
		t.Fatalf("ListFiles() failed, unexpected items count, expected: %d, actual: %d",limit,len(entries))
	}

	for _,entry:=range entries{
		t.Logf("ListItem:\n%s",entry.String())
	}
}
