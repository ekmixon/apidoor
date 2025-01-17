package managementapi_test

import (
	"bytes"
	"encoding/json"
	"github.com/future-architect/apidoor/managementapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestPostUser(t *testing.T) {
	if _, err := db.Exec("DELETE FROM apiuser"); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		contentType    string
		req            managementapi.PostUserReq
		wantHttpStatus int
		//wantRecord は期待されるDB作成レコードの値、idは比較対象外
		wantRecords []managementapi.User
	}{
		{
			name:        "正常に登録できる",
			contentType: "application/json",
			req: managementapi.PostUserReq{
				AccountID:    "user",
				EmailAddress: "test00@example.com",
				Password:     "password",
				Name:         "full name",
			},
			wantHttpStatus: http.StatusCreated,
			wantRecords: []managementapi.User{
				{
					AccountID:      "user",
					EmailAddress:   "test00@example.com",
					Name:           "full name",
					PermissionFlag: "00",
				},
			},
		},
		{
			name:        "パスワードに記号が含まれており、正常に登録できる",
			contentType: "application/json",
			req: managementapi.PostUserReq{
				AccountID:    "user1",
				EmailAddress: "test01@example.com",
				Password:     "p@ss12Word",
				Name:         "full name",
			},
			wantHttpStatus: http.StatusCreated,
			wantRecords: []managementapi.User{
				{
					AccountID:      "user1",
					EmailAddress:   "test01@example.com",
					Name:           "full name",
					PermissionFlag: "00",
				},
			},
		},
		{
			name:        "名前が設定されていなくても、正常に登録できる",
			contentType: "application/json",
			req: managementapi.PostUserReq{
				AccountID:    "user2",
				EmailAddress: "test02@example.com",
				Password:     "password",
				Name:         "",
			},
			wantHttpStatus: http.StatusCreated,
			wantRecords: []managementapi.User{
				{
					AccountID:      "user2",
					EmailAddress:   "test02@example.com",
					Name:           "",
					PermissionFlag: "00",
				},
			},
		},
		{
			name:        "必須項目に空欄があるとき、登録できない",
			contentType: "application/json",
			req: managementapi.PostUserReq{
				AccountID:    "",
				EmailAddress: "test03@example.com",
				Password:     "password",
				Name:         "full name",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantRecords:    []managementapi.User{},
		},
		{
			name:        "account_idにprintable ascii以外の文字が含まれていたとき、登録できない",
			contentType: "application/json",
			req: managementapi.PostUserReq{
				AccountID:    "userユーザー",
				EmailAddress: "test04@example.com",
				Password:     "password",
				Name:         "full name",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantRecords:    []managementapi.User{},
		},
		{
			name:        "email_addressの文字列がメールアドレスとして不正であるとき、登録できない",
			contentType: "application/json",
			req: managementapi.PostUserReq{
				AccountID:    "user",
				EmailAddress: "test05.@example.com",
				Password:     "password",
				Name:         "full name",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantRecords:    []managementapi.User{},
		},
		{
			name:        "Content-Typeがapplication/json以外であるとき、登録できない",
			contentType: "text/plain",
			req: managementapi.PostUserReq{
				AccountID:    "user",
				EmailAddress: "test06@example.com",
				Password:     "password",
				Name:         "full name",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantRecords:    []managementapi.User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tt.req)
			if err != nil {
				t.Errorf("create request body error: %v", err)
				return
			}
			body := bytes.NewReader(bodyBytes)

			r := httptest.NewRequest(http.MethodPost, "localhost:3000/mgmt/product", body)
			r.Header.Add("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			managementapi.PostUser(w, r)

			rw := w.Result()
			if rw.StatusCode != tt.wantHttpStatus {
				t.Errorf("wrong http status code: got %d, want %d", rw.StatusCode, tt.wantHttpStatus)
			}

			if rw.StatusCode != http.StatusCreated {
				return
			}

			rows, err := db.Queryx("SELECT * from apiuser WHERE email_address=$1", tt.req.EmailAddress)
			if err != nil {
				t.Errorf("db get products error: %v", err)
				return
			}

			list := []managementapi.User{}
			for rows.Next() {
				var row managementapi.User

				if err := rows.StructScan(&row); err != nil {
					t.Errorf("reading row error: %v", err)
					return
				}

				list = append(list, row)
			}

			if diff := cmp.Diff(tt.wantRecords, list,
				cmpopts.IgnoreFields(managementapi.User{}, "ID", "LoginPasswordHash",
					"CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("db get users responce differs:\n %v", diff)
			}

			// checking that passwords are stored in a hash
			hashRegex := regexp.MustCompile("\\$2a\\$\\w+\\$[\\w.]+")
			for _, v := range list {
				if !hashRegex.Match([]byte(v.LoginPasswordHash)) {
					t.Errorf("password hash format is wrong, got: %s", v.LoginPasswordHash)
				}
			}

		})
	}

	if _, err := db.Exec("DELETE FROM apiuser"); err != nil {
		t.Fatal(err)
	}

}
