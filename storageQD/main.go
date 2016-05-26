package main

import (
	"fmt"
    "google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
	"io"
    "google.golang.org/appengine"
    "golang.org/x/net/context"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

const gcsBucket = "scotty"

type demo struct {
	ctx    context.Context
	res    http.ResponseWriter
	bucket *storage.BucketHandle
	client *storage.Client
}

func handler(res http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	ctx := appengine.NewContext(req)

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "ERROR handler NewClient: ", err)
		return
	}
	defer client.Close()

	d := &demo{
		ctx:    ctx,
		res:    res,
		client: client,
		bucket: client.Bucket(gcsBucket),
	}

	d.delFiles()
	d.createFiles()
	d.listFiles()
	io.WriteString(d.res, "\nFILE NAMES WITH DELIMITER QUERY ( Delimeter: / )\n")
	d.listDelim()

}

func (d *demo) createFile(fileName string) {

	wc := d.bucket.Object(fileName).NewWriter(d.ctx)
	wc.ContentType = "text/plain"

	if _, err := wc.Write([]byte("abcde\n")); err != nil {
		log.Errorf(d.ctx, "unable to write data to bucket", gcsBucket, fileName, err)
		return
	}
	if err := wc.Close(); err != nil {
		log.Errorf(d.ctx, " unable to close bucket", gcsBucket, fileName, err)
		return
	}
}

func (d *demo) listDelim() {

	query := &storage.Query{
		Delimiter: "/",
	}

	objs, err := d.bucket.List(d.ctx, query)
	if err != nil {
		log.Errorf(d.ctx, "listBucketDirMode: unable to list bucket %q: %v", gcsBucket, err)
		return
	}

	for _, obj := range objs.Results {
		fmt.Fprintf(d.res, "%v\n", obj.Name)
	}

	fmt.Fprintf(d.res, "\nPREFIXES ( storage.ObjectList Prefixes )\n%v", objs.Prefixes)
}



func (d *demo) listFiles() {
	io.WriteString(d.res, "ALL FILE NAMES\n")

	objs, err := d.bucket.List(d.ctx, nil)
	if err != nil {
		log.Errorf(d.ctx, "%v", err)
		return
	}

	for _, obj := range objs.Results {
		io.WriteString(d.res, obj.Name+"\n")
	}
}

func (d *demo) createFiles() {
	for _, n := range []string{"firstobject", "folder1/folder2/secondobject",} {
		d.createFile(n)
	}
}

func (d *demo) delFiles() {
	objs, err := d.bucket.List(d.ctx, nil)
	if err != nil {
		log.Errorf(d.ctx, "%v", err)
		return
	}

	for _, obj := range objs.Results {
		if err := d.bucket.Object(obj.Name).Delete(d.ctx); err != nil {
			log.Errorf(d.ctx, "unable to delete bucket", d.bucket, obj.Name, err)
			return
		}
	}
}
