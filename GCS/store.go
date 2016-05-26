package main

import (
	"golang.org/x/net/context"
	"google.golang.org/cloud/storage"
	"io"
)

func keep(ctx context.Context, fname string, rdr io.Reader) error {
	client, err := storage.NewClient(ctx)
    if err != nil {
        return err
    }
    defer client.Close()

    writer := client.Bucket(bucket).Object(fname).NewWriter(ctx)
    writer.ACL = []storage.ACLRule{
        {storage.AllUsers, storage.RoleReader},
    }
    io.Copy(writer, rdr)
    return writer.Close()
}


func get(ctx context.Context, fname string) (io.ReadCloser, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return client.Bucket(bucket).Object(fname).NewReader(ctx)
}
