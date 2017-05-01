//
// exec_test.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This file is part of Roadie queue manager.
//
// Roadie Queue Manager is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Roadie Queue Manager is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Roadie queue manager. If not, see <http://www.gnu.org/licenses/>.
//

package main

import "testing"

func TestParseURL(t *testing.T) {
	var opt DownloadOpt

	// Basic URL
	opt = parseURL("http://www.sample.com/sample.txt")
	if opt.Src != "http://www.sample.com/sample.txt" || opt.Dest != "sample.txt" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// URL with renaming
	opt = parseURL("http://www.sample.com/sample.txt:another.txt")
	if opt.Src != "http://www.sample.com/sample.txt" || opt.Dest != "another.txt" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// URL with a destination folder
	opt = parseURL("http://www.sample.com/sample.txt:/tmp/")
	if opt.Src != "http://www.sample.com/sample.txt" || opt.Dest != "/tmp/sample.txt" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// URL with renaming and a destination folder
	opt = parseURL("http://www.sample.com/sample.txt:/tmp/another.txt")
	if opt.Src != "http://www.sample.com/sample.txt" || opt.Dest != "/tmp/another.txt" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Basic URL of a compressed file
	opt = parseURL("http://www.sample.com/sample.zip")
	if opt.Src != "http://www.sample.com/sample.zip" || opt.Dest != "sample.zip" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || !opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// URL of a compressed file with renaming
	opt = parseURL("http://www.sample.com/sample.zip:another.zip")
	if opt.Src != "http://www.sample.com/sample.zip" || opt.Dest != "another.zip" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// URL of a compressed file with a destination folder
	opt = parseURL("http://www.sample.com/sample.zip:/tmp/")
	if opt.Src != "http://www.sample.com/sample.zip" || opt.Dest != "/tmp/sample.zip" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || !opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// URL of a compressed file with renaming and a destination folder
	opt = parseURL("http://www.sample.com/sample.zip:/tmp/another.zip")
	if opt.Src != "http://www.sample.com/sample.zip" || opt.Dest != "/tmp/another.zip" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Basic URL of a tar+gzipped file
	opt = parseURL("http://www.sample.com/sample.tar.gz")
	if opt.Src != "http://www.sample.com/sample.tar.gz" || opt.Dest != "sample.tar.gz" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || !opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Basic URL of a tar file
	opt = parseURL("http://www.sample.com/sample.tar")
	if opt.Src != "http://www.sample.com/sample.tar" || opt.Dest != "sample.tar" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if !opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Dropbox URL with the host name.
	opt = parseURL("dropbox://www.dropbox.com/s/aaaaaaaaaa/sample.txt?dl=0")
	if opt.Src != "https://www.dropbox.com/s/aaaaaaaaaa/sample.txt?dl=1" || opt.Dest != "sample.txt" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Dropbox URL without the host name.
	opt = parseURL("dropbox://s/aaaaaaaaaa/sample.txt?dl=0")
	if opt.Src != "https://www.dropbox.com/s/aaaaaaaaaa/sample.txt?dl=1" || opt.Dest != "sample.txt" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Dropbox URL with renaming.
	opt = parseURL("dropbox://s/aaaaaaaaaa/sample.txt?dl=0:another.txt")
	if opt.Src != "https://www.dropbox.com/s/aaaaaaaaaa/sample.txt?dl=1" || opt.Dest != "another.txt" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Dropbox URL with a destination directory.
	opt = parseURL("dropbox://s/aaaaaaaaaa/sample.txt?dl=0:/tmp/")
	if opt.Src != "https://www.dropbox.com/s/aaaaaaaaaa/sample.txt?dl=1" || opt.Dest != "/tmp/sample.txt" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Dropbox URL with renaming and a destination directory.
	opt = parseURL("dropbox://s/aaaaaaaaaa/sample.txt?dl=0:/tmp/another.txt")
	if opt.Src != "https://www.dropbox.com/s/aaaaaaaaaa/sample.txt?dl=1" || opt.Dest != "/tmp/another.txt" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Dropbox folder URL with the host name.
	opt = parseURL("dropbox://www.dropbox.com/sh/aaaaaaaaaa/aaaaaaaaa?dl=0")
	if opt.Src != "https://www.dropbox.com/sh/aaaaaaaaaa/aaaaaaaaa?dl=1" || opt.Dest != "dropbox.zip" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || !opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Dropbox folder URL without the host name.
	opt = parseURL("dropbox://sh/aaaaaaaaaa/aaaaaaaaa?dl=0")
	if opt.Src != "https://www.dropbox.com/sh/aaaaaaaaaa/aaaaaaaaa?dl=1" || opt.Dest != "dropbox.zip" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || !opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Dropbox folder URL with renaming.
	opt = parseURL("dropbox://sh/aaaaaaaaaa/aaaaaaaaa?dl=0:another.zip")
	if opt.Src != "https://www.dropbox.com/sh/aaaaaaaaaa/aaaaaaaaa?dl=1" || opt.Dest != "another.zip" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Dropbox folder URL with a destination folder.
	opt = parseURL("dropbox://sh/aaaaaaaaaa/aaaaaaaaa?dl=0:/tmp/")
	if opt.Src != "https://www.dropbox.com/sh/aaaaaaaaaa/aaaaaaaaa?dl=1" || opt.Dest != "/tmp/dropbox.zip" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || !opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

	// Dropbox folder URL with renaming and a destination folder.
	opt = parseURL("dropbox://sh/aaaaaaaaaa/aaaaaaaaa?dl=0:/tmp/another.zip")
	if opt.Src != "https://www.dropbox.com/sh/aaaaaaaaaa/aaaaaaaaa?dl=1" || opt.Dest != "/tmp/another.zip" {
		t.Errorf("Parsed URL is not correct (src: %v, dest: %v)", opt.Src, opt.Dest)
	}
	if opt.Tar || opt.TarGz || opt.Zip {
		t.Errorf("Decompress configuration is not correct (dest:%v, tar: %v, tar.gz: %v, zip: %v)", opt.Dest, opt.Tar, opt.TarGz, opt.Zip)
	}

}
