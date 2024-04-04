package harbor

import (
	"encoding/json"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func SyncLabels(src, dst *Harbor, logger log.Logger) error {

	srcLabels, err := getLabels(src)
	if err != nil {
		return fmt.Errorf("cannot fetch %s labels: %s", *src.Host, err)
	}

	dstLabels, err := getLabels(dst)
	if err != nil {
		return fmt.Errorf("cannot fetch %s labels: %s", *dst.Host, err)
	}

	for labelID, label := range srcLabels {

		if _, ok := dstLabels[labelID]; ok {
			err = updateLabel(dst, label)
			if err != nil {
				return fmt.Errorf("unable to update label %s: %s", label.Name, err)
			}
			level.Debug(logger).Log("msg", fmt.Sprintf("Updated label %s to %s", label.Name, *dst.Host))
		} else {
			err = addLabel(dst, label)
			if err != nil {
				return fmt.Errorf("unable to add label %s: %s", label.Name, err)
			}
			level.Debug(logger).Log("msg", fmt.Sprintf("Added label %s to %s", label.Name, *dst.Host))
		}

	}

	if len(srcLabels) < len(dstLabels) {
		for labelID, label := range dstLabels {
			if _, ok := srcLabels[labelID]; !ok {
				err = deleteLabel(dst, label)
				if err != nil {
					return fmt.Errorf("unable to delete label %s: %s", label.Name, err)
				}
				level.Debug(logger).Log("msg", fmt.Sprintf("Deleted label %s from %s", label.Name, *dst.Host))
			}
		}
	}

	return nil
}

func getLabels(target *Harbor) (map[int]Label, error) {

	var (
		labels    map[int]Label
		rawLabels []Label
	)
	labels = make(map[int]Label)

	url := fmt.Sprintf("%s://%s/api/v2.0/labels?page=0&page_size=0&scope=g", *hbSchema, *target.Host)

	statusCode, body, err := request("GET", url, *target.User, *target.Pass, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching labels: %s", err)
	}

	if statusCode > 399 {
		return nil, fmt.Errorf("fetch labels returned %d", statusCode)
	}

	err = json.Unmarshal(*body, &rawLabels)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling labels: %s", err)
	}

	for _, label := range rawLabels {
		labels[label.ID] = label
	}

	return labels, nil

}

func updateLabel(dst *Harbor, label Label) error {

	url := fmt.Sprintf("%s://%s/api/v2.0/labels/%d", *hbSchema, *dst.Host, label.ID)

	rbody, err := json.Marshal(label)
	if err != nil {
		return fmt.Errorf("error marshalling label body: %s", err)
	}

	statusCode, _, err := request("PUT", url, *dst.User, *dst.Pass, rbody)
	if err != nil {
		return fmt.Errorf("error updating label %s: %s", label.Name, err)
	}

	if statusCode > 399 {
		return fmt.Errorf("update label %s returned %d", label.Name, statusCode)
	}

	return nil
}

func addLabel(dst *Harbor, label Label) error {

	url := fmt.Sprintf("%s://%s/api/v2.0/labels", *hbSchema, *dst.Host)

	rbody, err := json.Marshal(label)
	if err != nil {
		return fmt.Errorf("error marshalling label body: %s", err)
	}

	statusCode, _, err := request("POST", url, *dst.User, *dst.Pass, rbody)
	if err != nil {
		return fmt.Errorf("error updating label %s: %s", label.Name, err)
	}

	if statusCode > 399 {
		return fmt.Errorf("update label %s returned %d", label.Name, statusCode)
	}

	return nil
}

func deleteLabel(dst *Harbor, label Label) error {

	url := fmt.Sprintf("%s://%s/api/v2.0/labels/%d", *hbSchema, *dst.Host, label.ID)

	statusCode, _, err := request("DELETE", url, *dst.User, *dst.Pass, nil)
	if err != nil {
		return fmt.Errorf("error deleting label %s: %s", label.Name, err)
	}

	if statusCode > 399 {
		return fmt.Errorf("delete label %s returned %d", label.Name, statusCode)
	}

	return nil
}
