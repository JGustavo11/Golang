package main

import (
	"bytes"
	"fmt"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
	"golang.org/x/net/context"
	"io"
    "net/http"
	"strings"
    "google.golang.org/appengine/file"
	"io/ioutil"
	"google.golang.org/appengine"
	
)

var bucket = "bucky"

func init() {
	http.HandleFunc("/", handler)
}

type demo struct {
	bucket *storage.BucketHandle
	client *storage.Client

	w   http.ResponseWriter
	ctx context.Context
	
	cleanUp []string
	
	failed bool
}

func (d *demo) errorf(format string, args ...interface{}) {
	d.failed = true
	log.Errorf(d.ctx, format, args...)
}


func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ctx := appengine.NewContext(r)
	if bucket == "" {
		var err error
		if bucket, err = file.DefaultBucketName(ctx); err != nil {
			log.Errorf(ctx, "error", err)
			return
		}
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "error", err)
		return
	}
	defer client.Close()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Demo GCS Application running from Version: %v\n", appengine.VersionID(ctx))
	fmt.Fprintf(w, "Using bucket name: %v\n\n", bucket)

	d := &demo{
		w:      w,
		ctx:    ctx,
		client: client,
		bucket: client.Bucket(bucket),
	}

	n := "demo-testfile-go"
	d.makeF(n)
	d.rFile(n)
	d.cFile(n)
	d.statFile(n)
	d.createListFiles()
	d.listBucket()
	d.listBucketDirMode()
	d.defaultACL()
	d.putDefaultACLRule()
	d.deleteDefaultACLRule()
	d.bucketACL()
	d.putBucketACLRule()
	d.deleteBucketACLRule()
	d.acl(n)
	d.putACLRule(n)
	d.deleteACLRule(n)
	d.deleteFiles()

	if d.failed {
		io.WriteString(w, "\nDemo failed.\n")
	} else {
		io.WriteString(w, "\nDemo succeeded.\n")
	}
}

func (d *demo) makeF(fileName string) {
	fmt.Fprintf(d.w, "Creating file /%v/%v\n", bucket, fileName)

	wc := d.bucket.Object(fileName).NewWriter(d.ctx)
	wc.ContentType = "text/plain"
	wc.Metadata = map[string]string{
		"x-goog-meta-foo": "foo",
		"x-goog-meta-bar": "bar",
	}
	d.cleanUp = append(d.cleanUp, fileName)

	if _, err := wc.Write([]byte("abcde\n")); err != nil {
		d.errorf("makeF: unable to write data to bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	if _, err := wc.Write([]byte(strings.Repeat("f", 1024*4) + "\n")); err != nil {
		d.errorf("makeF: unable to write data to bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	if err := wc.Close(); err != nil {
		d.errorf("makeF: unable to close bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
}

func (d *demo) rFile(fileName string) {
	io.WriteString(d.w, "test")

	rc, err := d.bucket.Object(fileName).NewReader(d.ctx)
	if err != nil {
		d.errorf("cant read the file", bucket, fileName, err)
		return
	}
	defer rc.Close()
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		d.errorf("data in file is un readable", bucket, fileName, err)
		return
	}

	fmt.Fprintf(d.w, "%s\n", bytes.SplitN(slurp, []byte("\n"), 2)[0])
	if len(slurp) > 1024 {
		fmt.Fprintf(d.w, "...%s\n", slurp[len(slurp)-1024:])
	} else {
		fmt.Fprintf(d.w, "%s\n", slurp)
	}
}

func (d *demo) cFile(fileName string) {
	copyName := fileName + "-copy"
	fmt.Fprintf(d.w, "copying initiating..", bucket, fileName, bucket, copyName)

	obj, err := d.bucket.Object(fileName).CopyTo(d.ctx, d.bucket.Object(copyName), nil)
	if err != nil {
		d.errorf("Cannot copy files", bucket, fileName, bucket, copyName, err)
		return
	}
	d.cleanUp = append(d.cleanUp, copyName)

	d.dumpStats(obj)
}

func (d *demo) dumpStats(obj *storage.ObjectAttrs) {
	fmt.Fprintf(d.w, "(filename: /%v/%v, ", obj.Bucket, obj.Name)
	fmt.Fprintf(d.w, "ContentType: %q, ", obj.ContentType)
	fmt.Fprintf(d.w, "ACL: %#v, ", obj.ACL)
	fmt.Fprintf(d.w, "Owner: %v, ", obj.Owner)
	fmt.Fprintf(d.w, "ContentEncoding: %q, ", obj.ContentEncoding)
	fmt.Fprintf(d.w, "Size: %v, ", obj.Size)
	fmt.Fprintf(d.w, "MD5: %q, ", obj.MD5)
	fmt.Fprintf(d.w, "CRC32C: %q, ", obj.CRC32C)
	fmt.Fprintf(d.w, "Metadata: %#v, ", obj.Metadata)
	fmt.Fprintf(d.w, "MediaLink: %q, ", obj.MediaLink)
	fmt.Fprintf(d.w, "StorageClass: %q, ", obj.StorageClass)
	if !obj.Deleted.IsZero() {
		fmt.Fprintf(d.w, "Deleted: %v, ", obj.Deleted)
	}
	fmt.Fprintf(d.w, "Updated: %v)\n", obj.Updated)
}

func (d *demo) statFile(fileName string) {
	io.WriteString(d.w, "\nFile stat:\n")

	obj, err := d.bucket.Object(fileName).Attrs(d.ctx)
	if err != nil {
		d.errorf("statFile: unable to stat file from bucket %q, file %q: %v", bucket, fileName, err)
		return
	}

	d.dumpStats(obj)
}

func (d *demo) createListFiles() {
	io.WriteString(d.w, "\nCreating more files for listbucket...\n")
	for _, n := range []string{"foo1", "foo2", "bar", "bar/1", "bar/2", "boo/"} {
		d.makeF(n)
	}
}

func (d *demo) listBucket() {
	io.WriteString(d.w, "\nListbucket result:\n")

	query := &storage.Query{Prefix: "foo"}
	for query != nil {
		objs, err := d.bucket.List(d.ctx, query)
		if err != nil {
			d.errorf("listBucket: unable to list bucket %q: %v", bucket, err)
			return
		}
		query = objs.Next

		for _, obj := range objs.Results {
			d.dumpStats(obj)
		}
	}
}

func (d *demo) listDir(name, indent string) {
	query := &storage.Query{Prefix: name, Delimiter: "/"}
	for query != nil {
		objs, err := d.bucket.List(d.ctx, query)
		if err != nil {
			d.errorf("listBucketDirMode: unable to list bucket %q: %v", bucket, err)
			return
		}
		query = objs.Next

		for _, obj := range objs.Results {
			fmt.Fprint(d.w, indent)
			d.dumpStats(obj)
		}
		for _, dir := range objs.Prefixes {
			fmt.Fprintf(d.w, "%v(directory: /%v/%v)\n", indent, bucket, dir)
			d.listDir(dir, indent+"  ")
		}
	}
}

func (d *demo) listBucketDirMode() {
	io.WriteString(d.w, "\nListbucket directory mode result:\n")
	d.listDir("b", "")
}

func (d *demo) dumpDefaultACL() {
	acl, err := d.bucket.ACL().List(d.ctx)
	if err != nil {
		d.errorf("defaultACL: unable to list default object ACL for bucket %q: %v", bucket, err)
		return
	}
	for _, v := range acl {
		fmt.Fprintf(d.w, "Scope: %q, Permission: %q\n", v.Entity, v.Role)
	}
}

func (d *demo) defaultACL() {
	io.WriteString(d.w, "\nDefault object ACL:\n")
	d.dumpDefaultACL()
}

for this bucket.
func (d *demo) putDefaultACLRule() {
	io.WriteString(d.w, "\nPut Default object ACL Rule:\n")
	err := d.bucket.DefaultObjectACL().Set(d.ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		d.errorf("putDefaultACLRule: unable to save default object ACL rule for bucket %q: %v", bucket, err)
		return
	}
	d.dumpDefaultACL()
}

func (d *demo) deleteDefaultACLRule() {
	io.WriteString(d.w, "\nDelete Default object ACL Rule:\n")
	err := d.bucket.DefaultObjectACL().Delete(d.ctx, storage.AllUsers)
	if err != nil {
		d.errorf("deleteDefaultACLRule: unable to delete default object ACL rule for bucket %q: %v", bucket, err)
		return
	}
	d.dumpDefaultACL()
}

func (d *demo) dumpBucketACL() {
	acl, err := d.bucket.ACL().List(d.ctx)
	if err != nil {
		d.errorf("dumpBucketACL: unable to list bucket ACL for bucket %q: %v", bucket, err)
		return
	}
	for _, v := range acl {
		fmt.Fprintf(d.w, "Scope: %q, Permission: %q\n", v.Entity, v.Role)
	}
}

func (d *demo) bucketACL() {
	io.WriteString(d.w, "\nBucket ACL:\n")
	d.dumpBucketACL()
}

func (d *demo) putBucketACLRule() {
	io.WriteString(d.w, "\nPut Bucket ACL Rule:\n")
	err := d.bucket.ACL().Set(d.ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		d.errorf("putBucketACLRule: unable to save bucket ACL rule for bucket %q: %v", bucket, err)
		return
	}
	d.dumpBucketACL()
}

func (d *demo) deleteBucketACLRule() {
	io.WriteString(d.w, "\nDelete Bucket ACL Rule:\n")
	err := d.bucket.ACL().Delete(d.ctx, storage.AllUsers)
	if err != nil {
		d.errorf("deleteBucketACLRule: unable to delete bucket ACL rule for bucket %q: %v", bucket, err)
		return
	}
	d.dumpBucketACL()
}

func (d *demo) dumpACL(fileName string) {
	acl, err := d.bucket.Object(fileName).ACL().List(d.ctx)
	if err != nil {
		d.errorf("dumpACL: unable to list file ACL for bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	for _, v := range acl {
		fmt.Fprintf(d.w, "Scope: %q, Permission: %q\n", v.Entity, v.Role)
	}
}

func (d *demo) acl(fileName string) {
	fmt.Fprintf(d.w, "\nACL for file %v:\n", fileName)
	d.dumpACL(fileName)
}

func (d *demo) putACLRule(fileName string) {
	fmt.Fprintf(d.w, "\nPut ACL rule for file %v:\n", fileName)
	err := d.bucket.Object(fileName).ACL().Set(d.ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		d.errorf("putACLRule: unable to save ACL rule for bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	d.dumpACL(fileName)
}

func (d *demo) deleteACLRule(fileName string) {
	fmt.Fprintf(d.w, "\nDelete ACL rule for file %v:\n", fileName)
	err := d.bucket.Object(fileName).ACL().Delete(d.ctx, storage.AllUsers)
	if err != nil {
		d.errorf("deleteACLRule: unable to delete ACL rule for bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	d.dumpACL(fileName)
}
func (d *demo) deleteFiles() {
	io.WriteString(d.w, "\nDeleting files...\n")
	for _, v := range d.cleanUp {
		fmt.Fprintf(d.w, "Deleting file %v\n", v)
		if err := d.bucket.Object(v).Delete(d.ctx); err != nil {
			d.errorf("deletion cannot be completed", bucket, v, err)
			return
		}
	}
}
