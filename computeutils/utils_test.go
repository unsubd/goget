package computeutils

import "testing"

func TestFileNameFromUrl(t *testing.T) {
	url := "file:///Users/adit/workspace/localhost.htm"
	name := FileNameFromUrl(url)

	if name != "localhost.htm" {
		t.Fail()
	}

	url = "http://localhost:8000/hello/abc/"
	name = FileNameFromUrl(url)
	if name != "abc" {
		t.Fail()
	}
}

func TestCreateBatches(t *testing.T) {
	batches := CreateBatches(129, 7)
	if len(batches) != 17 && batches[len(batches)-1][0] == 128 && batches[len(batches)-1][1] == 135 {
		t.Fail()
	}
}

func TestGetFilePath(t *testing.T) {
	name := GetFilePath("dir", "file")
	if name != "dir/file" {
		t.Error("Expected", "dir/file", "Returned", name)
	}

	name = GetFilePath("/dir", "file")
	if name != "/dir/file" {
		t.Error("Expected", "dir/file", "Returned", name)
	}

	name = GetFilePath("/dir/", "file")
	if name != "/dir/file" {
		t.Error("Expected", "dir/file", "Returned", name)
	}
	name = GetFilePath("/dir/", "/file")
	if name != "/dir/file" {
		t.Error("Expected", "dir/file", "Returned", name)
	}
}
