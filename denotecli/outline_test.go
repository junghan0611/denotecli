// outline_test.go
package main

import (
	"testing"
)

func TestExtractOutline(t *testing.T) {
	content := `#+title: 다단계 문서 테스트
#+filetags: :test:outline:
#+identifier: 20260101T120000

* 1장 서론
서론 본문
** 1.1 배경
배경 설명
** 1.2 목적 :IMPORTANT:
목적 설명
* 2장 본론
본론 시작
** 2.1 방법론
방법론 설명
*** 2.1.1 세부 방법
세부 내용
** 2.2 결과
결과 설명
* 3장 결론
결론 본문`

	outline := ExtractOutline(content)

	if len(outline) != 8 {
		t.Fatalf("expected 8 headings, got %d", len(outline))
	}

	// Check levels
	expectedLevels := []int{1, 2, 2, 1, 2, 3, 2, 1}
	for i, e := range outline {
		if e.Level != expectedLevels[i] {
			t.Errorf("outline[%d].Level = %d, want %d", i, e.Level, expectedLevels[i])
		}
	}

	// Check first heading
	if outline[0].Title != "1장 서론" {
		t.Errorf("outline[0].Title = %q, want %q", outline[0].Title, "1장 서론")
	}
	if outline[0].Line != 5 {
		t.Errorf("outline[0].Line = %d, want 5", outline[0].Line)
	}

	// Check heading with org tags
	if outline[2].Title != "1.2 목적" {
		t.Errorf("outline[2].Title = %q, want %q", outline[2].Title, "1.2 목적")
	}
	if outline[2].Tags != ":IMPORTANT:" {
		t.Errorf("outline[2].Tags = %q, want %q", outline[2].Tags, ":IMPORTANT:")
	}

	// Check deepest heading
	if outline[5].Level != 3 {
		t.Errorf("outline[5].Level = %d, want 3", outline[5].Level)
	}
	if outline[5].Title != "2.1.1 세부 방법" {
		t.Errorf("outline[5].Title = %q, want %q", outline[5].Title, "2.1.1 세부 방법")
	}

	// Last heading
	if outline[7].Title != "3장 결론" {
		t.Errorf("outline[7].Title = %q, want %q", outline[7].Title, "3장 결론")
	}
}

func TestExtractOutlineEmpty(t *testing.T) {
	content := `#+title: 빈 문서
#+filetags: :test:
#+identifier: 20260101T130000

프론트매터만 있고 헤딩 없는 문서.`

	outline := ExtractOutline(content)
	if len(outline) != 0 {
		t.Errorf("expected 0 headings, got %d", len(outline))
	}
}

func TestExtractOutlineSingleHeading(t *testing.T) {
	content := "* 유일한 헤딩\n본문\n"
	outline := ExtractOutline(content)
	if len(outline) != 1 {
		t.Fatalf("expected 1 heading, got %d", len(outline))
	}
	if outline[0].Title != "유일한 헤딩" {
		t.Errorf("Title = %q", outline[0].Title)
	}
	if outline[0].Level != 1 {
		t.Errorf("Level = %d", outline[0].Level)
	}
	if outline[0].Line != 1 {
		t.Errorf("Line = %d", outline[0].Line)
	}
}

func TestReadOutlineByID(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	do, err := ReadOutlineByID(files, "20260101T120000")
	if err != nil {
		t.Fatal(err)
	}
	if do.Title != "다단계 문서 테스트" {
		t.Errorf("Title = %q", do.Title)
	}
	if len(do.Outline) != 8 {
		t.Errorf("expected 8 headings, got %d", len(do.Outline))
	}
	// Check link extraction works too
	if len(do.Links) != 1 || do.Links[0] != "20251107T082610" {
		t.Errorf("Links = %v, want [20251107T082610]", do.Links)
	}
}

func TestReadOutlineByIDNotFound(t *testing.T) {
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	_, err := ReadOutlineByID(files, "99999999T999999")
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}

func TestReadOutlineByIDSimple(t *testing.T) {
	// Test with a file that has single heading (에릭 호퍼)
	dir := setupTestDir(t)
	files := ScanDirs([]string{dir})

	do, err := ReadOutlineByID(files, "20251107T082610")
	if err != nil {
		t.Fatal(err)
	}
	if len(do.Outline) != 1 {
		t.Fatalf("expected 1 heading, got %d", len(do.Outline))
	}
	if do.Outline[0].Title != "본문" {
		t.Errorf("Title = %q, want 본문", do.Outline[0].Title)
	}
}

func TestExtractOutlineMultipleTags(t *testing.T) {
	content := "* TODO 작업 항목 :work:urgent:\n** DONE 완료 :done:\n"
	outline := ExtractOutline(content)
	if len(outline) != 2 {
		t.Fatalf("expected 2, got %d", len(outline))
	}
	if outline[0].Tags != ":work:urgent:" {
		t.Errorf("Tags = %q, want :work:urgent:", outline[0].Tags)
	}
	if outline[0].Title != "TODO 작업 항목" {
		t.Errorf("Title = %q", outline[0].Title)
	}
}
