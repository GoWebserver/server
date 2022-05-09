package settings

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"server/src"
	"server/src/log"
)

type Mime struct {
	Regex *regexp.Regexp
	Type  string
}

func LoadMimetypes() error {
	now := time.Now()

	//language=SQL
	sess := src.Session.Query(
		"SELECT regex, type, \"index\" FROM server.mime",
	)
	iter := sess.Iter()
	sett.Mimetypes.data = make([]Mime, iter.NumRows())
	for {
		row := make(map[string]any)
		if !iter.MapScan(row) {
			break
		}
		index, _ := strconv.Atoi(fmt.Sprintf("%d", row["index"]))
		sett.Mimetypes.data[index] = Mime{
			Regex: regexp.MustCompile(fmt.Sprintf("%s", row["regex"])),
			Type:  fmt.Sprintf("%s", row["type"]),
		}
	}
	if err := iter.Close(); err != nil {
		log.Err(err, "Error loading Mimetypes from DB")
		log.Debug(iter.Warnings())
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	log.Debug("Loaded Mimetypes in", time.Since(now))
	return nil
}

func LoadDefaultSite() error {
	now := time.Now()
	name := "DefaultSite" // IDE complains, it thinks "DefaultSite" is an SQL statement

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading DefaultSite from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	sett.DefaultSite.data = fmt.Sprintf("%s", setting["value"])

	log.Debug("Loaded DefaultSite in", time.Since(now))
	return nil
}

func LoadServerOff() error {
	now := time.Now()
	name := "ServerOff" // IDE complains, it thinks "DefaultSite" is an SQL statement

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading ServerOff from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	b, err := strconv.ParseBool(fmt.Sprintf("%s", setting["value"]))
	if err != nil {
		log.Err(err, "Invalid Value for ServerOff")
		return err
	}
	sett.ServerOff.data = b

	log.Debug("Loaded ServerOff in", time.Since(now))
	return nil
}

func LoadDeflateCompressMinSize() error {
	now := time.Now()
	name := "DeflateCompressMinSize" // IDE complains, it thinks "DefaultSite" is an SQL statement

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading DeflateCompressMinSize from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	b, err := strconv.ParseInt(fmt.Sprintf("%s", setting["value"]), 10, 64)
	if err != nil {
		log.Err(err, "Invalid Value for DeflateCompressMinSize")
		return err
	}
	sett.DeflateCompressMinSize.data = uint64(b)

	log.Debug("Loaded DeflateCompressMinSize in", time.Since(now))
	return nil
}

func LoadGZipCompressMinSize() error {
	now := time.Now()
	name := "GZipCompressMinSize" // IDE complains, it thinks "DefaultSite" is an SQL statement

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading GZipCompressMinSize from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	b, err := strconv.ParseInt(fmt.Sprintf("%s", setting["value"]), 10, 64)
	if err != nil {
		log.Err(err, "Invalid Value for GZipCompressMinSize")
		return err
	}
	sett.GZipCompressMinSize.data = uint64(b)

	log.Debug("Loaded GZipCompressMinSize in", time.Since(now))
	return nil
}

func LoadBrotliCompressMinSize() error {
	now := time.Now()
	name := "BrotliCompressMinSize" // IDE complains, it thinks "DefaultSite" is an SQL statement

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading BrotliCompressMinSize from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	b, err := strconv.ParseInt(fmt.Sprintf("%s", setting["value"]), 10, 64)
	if err != nil {
		log.Err(err, "Invalid Value for BrotliCompressMinSize")
		return err
	}
	sett.BrotliCompressMinSize.data = uint64(b)

	log.Debug("Loaded BrotliCompressMinSize in", time.Since(now))
	return nil
}

func LoadDeflateCompressMinCompression() error {
	now := time.Now()
	name := "DeflateCompressMinCompression" // IDE complains, it thinks "DefaultSite" is an SQL statement

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading DeflateCompressMinCompression from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	b, err := strconv.ParseFloat(fmt.Sprintf("%s", setting["value"]), 32)
	if err != nil {
		log.Err(err, "Invalid Value for DeflateCompressMinCompression")
		return err
	}
	sett.DeflateCompressMinCompression.data = float32(b)

	log.Debug("Loaded DeflateCompressMinCompression in", time.Since(now))
	return nil
}

func LoadGZipCompressMinCompression() error {
	now := time.Now()
	name := "GZipCompressMinCompression" // IDE complains, it thinks "DefaultSite" is an SQL statement

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading GZipCompressMinCompression from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	b, err := strconv.ParseFloat(fmt.Sprintf("%s", setting["value"]), 32)
	if err != nil {
		log.Err(err, "Invalid Value for GZipCompressMinCompression")
		return err
	}
	sett.GZipCompressMinCompression.data = float32(b)

	log.Debug("Loaded GZipCompressMinCompression in", time.Since(now))
	return nil
}

func LoadBrotliCompressMinCompression() error {
	now := time.Now()
	name := "BrotliCompressMinCompression" // IDE complains, it thinks "DefaultSite" is an SQL statement

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading BrotliCompressMinCompression from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	b, err := strconv.ParseFloat(fmt.Sprintf("%s", setting["value"]), 32)
	if err != nil {
		log.Err(err, "Invalid Value for BrotliCompressMinCompression")
		return err
	}
	sett.BrotliCompressMinCompression.data = float32(b)

	log.Debug("Loaded BrotliCompressMinCompression in", time.Since(now))
	return nil
}

func LoadEnableDeflateCompression() error {
	now := time.Now()
	name := "EnableDeflateCompression" // IDE complains, it thinks "DefaultSite" is an SQL statement

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading EnableDeflateCompression from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	b, err := strconv.ParseBool(fmt.Sprintf("%s", setting["value"]))
	if err != nil {
		log.Err(err, "Invalid Value for EnableDeflateCompression")
		return err
	}
	sett.EnableDeflateCompression.data = b

	log.Debug("Loaded EnableDeflateCompression in", time.Since(now))
	return nil
}

func LoadEnableGZipCompression() error {
	now := time.Now()
	name := "EnableGZipCompression" // IDE complains, it thinks "DefaultSite" is an SQL statement

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading EnableGZipCompression from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	b, err := strconv.ParseBool(fmt.Sprintf("%s", setting["value"]))
	if err != nil {
		log.Err(err, "Invalid Value for EnableGZipCompression")
		return err
	}
	sett.EnableGZipCompression.data = b

	log.Debug("Loaded EnableGZipCompression in", time.Since(now))
	return nil
}

func LoadEnableBrotliCompression() error {
	now := time.Now()
	name := "EnableBrotliCompression" // IDE complains, it thinks "DefaultSite" is an SQL statement

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading EnableBrotliCompression from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	b, err := strconv.ParseBool(fmt.Sprintf("%s", setting["value"]))
	if err != nil {
		log.Err(err, "Invalid Value for EnableBrotliCompression")
		return err
	}
	sett.EnableBrotliCompression.data = b

	log.Debug("Loaded EnableBrotliCompression in", time.Since(now))
	return nil
}
