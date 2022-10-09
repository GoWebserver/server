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

// Forbidden stores information about a rule to
// prevent access to a specific route
//
// FileExtension, AbsoluteFile and AbsoluteDirectory store
// their information inside Data, Regex stores inside Regex
type Forbidden struct {
	Data  string
	Regex *regexp.Regexp
	Type  ForbiddenType
}

type ForbiddenType string

const (
	FileExtension     ForbiddenType = "fe"
	Regex             ForbiddenType = "r"
	AbsoluteFile      ForbiddenType = "af"
	AbsoluteDirectory ForbiddenType = "ad"
)

func LoadMimetypes() error {
	now := time.Now()

	//language=SQL
	sess := src.Session.Query(
		"SELECT extension, mimetype, \"index\" FROM server.mime",
	)
	iter := sess.Iter()
	sett.Mimetypes.data = make([]Mime, iter.NumRows())
	for {
		row := map[string]any{}
		if !iter.MapScan(row) {
			break
		}
		index, _ := strconv.Atoi(fmt.Sprintf("%d", row["index"]))
		sett.Mimetypes.data[index] = Mime{
			Regex: regexp.MustCompile(fmt.Sprintf("%s", row["extension"])),
			Type:  fmt.Sprintf("%s", row["mimetype"]),
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
	name := "ServerOff"

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
	name := "DeflateCompressMinSize"

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
	name := "GZipCompressMinSize"

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
	name := "BrotliCompressMinSize"

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
	name := "DeflateCompressMinCompression"

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
	name := "GZipCompressMinCompression"

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
	name := "BrotliCompressMinCompression"

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
	name := "EnableDeflateCompression"

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
	name := "EnableGZipCompression"

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
	name := "EnableBrotliCompression"

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

func LoadForbidden() error {
	now := time.Now()

	//language=SQL
	sess := src.Session.Query(
		"SELECT \"index\", type, data FROM server. forbidden",
	)
	iter := sess.Iter()
	sett.Forbidden.data = make([]Forbidden, iter.NumRows())
	for {
		row := make(map[string]any)
		if !iter.MapScan(row) {
			break
		}
		index, _ := strconv.Atoi(fmt.Sprintf("%d", row["index"]))
		if row["type"] == Regex {
			sett.Forbidden.data[index] = Forbidden{
				Regex: regexp.MustCompile(fmt.Sprintf("%s", row["data"])),
				Type:  Regex,
			}
		} else {
			sett.Forbidden.data[index] = Forbidden{
				Data: fmt.Sprintf("%s", row["data"]),
				Type: ForbiddenType(fmt.Sprintf("%s", row["type"])),
			}
		}
	}
	if err := iter.Close(); err != nil {
		log.Err(err, "Error loading Forbidden from DB")
		log.Debug(iter.Warnings())
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	log.Debug("Loaded Forbidden in", time.Since(now))
	return nil
}

func LoadMaxURILength() error {
	now := time.Now()
	name := "MaxURILength"

	//language=SQL
	sess := src.Session.Query(
		"SELECT name, value FROM server.settings WHERE name=?", name,
	)
	setting := map[string]any{}
	err := sess.MapScan(setting)
	if err != nil {
		log.Err(err, "Error loading MaxURILength from DB")
		log.Debug(fmt.Sprintf("%s, attempts %d, latency: %dns", sess.String(), sess.Attempts(), sess.Latency()))
		return err
	}
	b, err := strconv.Atoi(fmt.Sprintf("%s", setting["value"]))
	if err != nil {
		log.Err(err, "Invalid Value for MaxURILength")
		return err
	}
	sett.MaxURILength.data = uint16(b)

	log.Debug("Loaded MaxURILength in", time.Since(now))
	return nil
}
