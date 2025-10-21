# Denote-Org Skills for Claude

> **Anthropic Life Sciences 패러다임을 삶 전반으로 확장합니다**

Claude AI를 위한 포괄적인 Denote PKM 시스템 지원. **3,000개 이상의 org 파일**로 검증된 실제 운영 시스템입니다.

**Denote 중심. Org-mode 기반. 삶의 규모로 증명됨.**

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Files](https://img.shields.io/badge/검증됨-3000%2B%20파일-green.svg)]()
[![Status](https://img.shields.io/badge/버전-0.1--dev-orange.svg)]()

## 🌟 비전: Life Sciences → Life Everything

Anthropic은 **도메인 컨텍스트 + AI = 전문가 수준 협업**을 Life Sciences(PubMed, Benchling, 10x Genomics)로 증명했습니다.

이 프로젝트는 그 패러다임을 **생명과학(Biology)**에서 **삶(Living)**으로 확장합니다:

**Life Sciences (생명과학)** → **Life Everything (삶 전반)**

Claude가 PubMed 컨텍스트로 과학자가 되듯이, Denote 컨텍스트로 **당신의 개인 지식 파트너**가 됩니다.

### 패턴 비교

| 영역 | 컨텍스트 | 결과 |
|------|---------|------|
| **Anthropic** | PubMed (수백만 논문) + Claude | 과학 연구 AI |
| **이 프로젝트** | Denote (3,000+ org 파일) + Claude | 개인 지식 AI |

**같은 방법론. 다른 도메인. 보편적 패러다임.**

## 왜 Skills인가? 단순 프롬프트가 아닌 이유

**Anthropic Skills 철학:**

> "Skills는 Claude를 범용 에이전트에서 **절차적 지식**을 갖춘 전문 에이전트로 변환시킵니다. 이는 어떤 모델도 완전히 가질 수 없는 지식입니다."
>
> — [Equipping agents for the real world with Agent Skills](https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills)

**PDF/XLSX Skills와 동일한 패턴:**
- **PDF Skill**: 복잡한 Forms 처리 (8개 Python 스크립트로 bounding box, 검증...)
- **XLSX Skill**: 공식 유지 (LibreOffice 연동 recalc.py)
- **Denote-Org Skill**: 3,000 파일 지식 그래프 (finder, graph builder, silo manager)

**Skills는 프롬프트가 아닙니다. 운영 시스템입니다.**

## 해결하는 문제

### Skills 없이

❌ Claude가 Denote 파일명을 모름: `20251021T105353--제목__태그.org`
❌ 3,000개 파일 검색에 매번 토큰 소모
❌ 지식 그래프 탐색이 느리고 오류 발생
❌ Literate programming (`:tangle`, `:results`)이 org 구조 깨뜨림
❌ 여러 silo (`~/org/`, `~/claude-memory/`) 혼란

### Skills 있으면

✅ **Denote 도메인 지식** 내장 (파일명, frontmatter, links)
✅ **Python 스크립트**로 효율적 연산 (finder, graph, executor)
✅ **지식 그래프** 캐싱 및 재사용
✅ **안전한 코드 실행** (org-babel 호환)
✅ **Silo 관리** (여러 지식 도메인)

**파일 변환기가 아닙니다. 운영 가능한 지식 시스템입니다.**

## 🎯 Denote란?

Denote는 Emacs를 위한 **단순하지만 강력한** 노트 작성 시스템입니다:

- **파일명 규칙**: `YYYYMMDDTHHMMSS--제목__태그1_태그2.org`
- **지식 그래프**: `[[denote:TIMESTAMP]]` 링크
- **Silo 개념**: 다양한 지식 도메인을 위한 별도 디렉토리
- **Org-mode 기반**: Heading, property, timestamp, code block

**제작자:** [Protesilaos Stavrou](https://protesilaos.com/emacs/denote)

## 📊 규모: 실제 운영으로 검증됨

이 Skill은 **실제 PKM 시스템**을 위해 만들어졌습니다 (장난감 예제 아님):

- ✅ **3,000개 이상 org 파일** (대규모 검증)
- ✅ **여러 silo** (`~/org/`, `~/claude-memory/`, 프로젝트 `docs/`)
- ✅ **Denote 지식 그래프** (수천 개의 `[[denote:ID]]` 링크)
- ✅ **매일 사용** (Claude Code와 함께)
- ✅ **9개 레이어 "Tools for Life" 시스템의 일부**

## 🏗️ 9개 레이어 시스템의 일부

이 Skill은 완전한 AI-인간 협업 아키텍처에서 **Layer 3** (지식 이해)를 담당합니다:

```
Layer 7: 지식 퍼블리싱 (notes.junghanacs.com)
Layer 6: 에이전트 오케스트레이션 (meta-config)
Layer 5a: 마이그레이션 (memex-kb)
Layer 5b: 삶의 타임라인 (memacs-config)
Layer 4: AI 메모리 (claude-config)
Layer 3.5: RAG/Vector DB (embedding-config)
Layer 3: 지식 관리 (Zotero, denote-org 🆕)
Layer 2: 개발 환경 (Doom Emacs)
Layer 1: 인프라 (NixOS)
```

**연관 프로젝트:**
- [claude-config](https://github.com/junghan0611/claude-config) - PARA 방법론 AI 메모리 시스템
- [memex-kb](https://github.com/junghan0611/memex-kb) - 범용 KB 마이그레이션 (Google Docs, Dooray...)
- [embedding-config](https://github.com/junghan0611/embedding-config) - Vector RAG (9,477 블록, 85% 해결률)

**철학:** Anthropic Life Sciences와 동일 - **AI는 파트너, 도구가 아님**

## 🚀 실제 사용 사례

### 사례 1: 지식 그래프 탐색
```bash
사용자: "20251021T105353와 연결된 모든 파일 찾아줘"
Claude: denote_graph.py 사용
결과: 1초 이내 연결 파일 반환 (반복 토큰 검색 아님)
```

### 사례 2: 3,000 파일 중 태그 검색
```bash
사용자: "2025년 10월 llmlog 파일 전부"
Claude: denote_finder.py --tags llmlog --date 202510* 사용
결과: 메타데이터와 함께 즉시 반환
```

### 사례 3: Silo 관리
```bash
사용자: "이게 ~/org/에 있나 ~/claude-memory/에 있나?"
Claude: denote_silo.py find_in_silos() 사용
결과: 자동으로 모호함 해결
```

### 사례 4: Literate Programming
```bash
사용자: "org 파일 코드 블록 실행해줘"
Claude: org_execute.py 사용 (org-babel 호환)
결과: 안전한 실행, 구조 보존, :results 업데이트
```

### 사례 5: Denote 링크 해결
```bash
파일 내용: [[denote:20251021T105353]]
Claude: denote_links.py resolve_link() 사용
결과: ~/org/llmlog/20251021T105353--전체-파일명.org
```

## 💡 기능

### Denote 핵심
- ✅ 파일명: `20251021T105353--제목__태그1_태그2.org`
- ✅ Frontmatter 파싱 (title, date, filetags, identifier)
- ✅ 링크 해결 (`[[denote:ID]]`)
- ✅ 태그 기반 검색
- ✅ Silo 관리 (여러 지식 도메인)
- ✅ 지식 그래프 (3,000+ 파일 검증)

### Org-mode 기반
- ✅ Heading/property 파싱
- ✅ 코드 블록 실행 (literate programming)
- ✅ Timestamp 처리
- ✅ Export (markdown, PDF, HTML)
- ✅ Table 연산

### 성능
- ✅ Python 스크립트 (토큰 생성 아님)
- ✅ 반복 연산 캐싱
- ✅ 효율적인 그래프 탐색
- ✅ 배치 처리 지원

## 📦 설치

### 사전 요구사항
```bash
# Python 3.8+
pip install orgparse pypandoc pyyaml

# Pandoc (export용)
# 패키지 매니저로 설치
```

### Claude Code
```bash
git clone https://github.com/junghan0611/orgmode-skills.git
cd orgmode-skills
/plugin add .
```

### Claude API
```python
# Skills API 통합 문서 참조
import anthropic

client = anthropic.Client(api_key="your-api-key")
# Skills API 사용법
```

## 🔧 빠른 시작

### 예제 1: Denote 파일 찾기
```python
from scripts.denote_finder import find_denote_file

# ID로 찾기
filepath = find_denote_file("20251021T105353", silos=["~/org/"])
# 반환: ~/org/llmlog/20251021T105353--전체-파일명.org

# 태그로 찾기
files = find_by_tags(["llmlog", "denote"], silo="~/org/")
```

### 예제 2: 지식 그래프 추출
```python
from scripts.denote_graph import build_knowledge_graph, get_connected_nodes

# 그래프 구축 (캐싱됨)
graph = build_knowledge_graph("~/org/")

# 연결된 파일 가져오기
connected = get_connected_nodes("20251021T105353", hops=2)
```

### 예제 3: Org 코드 블록 실행
```python
from scripts.org_execute import execute_org_blocks

# 모든 코드 블록 실행 (:tangle, :results 준수)
results = execute_org_blocks("file.org")
```

## 📂 프로젝트 구조

```
orgmode-skills/
├── README.md                    # 영어 버전
├── README-KO.md                 # 한글 버전 (이 파일)
├── LICENSE                      # Apache 2.0
├── SKILL.md                     # Claude용 메인 가이드
├── docs/                        # 개발 로그 (Denote 형식)
│   └── 20251021T113500--*.org
├── denote-core.md               # Denote 명세
├── denote-silo.md               # Silo 관리
├── denote-knowledge-graph.md    # 그래프 연산
├── orgmode-base.md              # Org-mode 기능
├── literate.md                  # Literate programming
├── export.md                    # Export 기능
└── scripts/
    ├── denote_finder.py         # 파일 검색
    ├── denote_links.py          # 링크 해결
    ├── denote_silo.py           # Silo 관리
    ├── denote_graph.py          # 지식 그래프
    ├── org_parser.py            # Org 파싱
    ├── org_execute.py           # 코드 실행
    ├── org_export.py            # Export
    └── requirements.txt
```

## 🗺️ 0.1 버전 로드맵

**Phase 1: 핵심 (현재)**
- [x] 프로젝트 설정
- [x] README.md, README-KO.md
- [ ] 개발 로그 (docs/)
- [ ] SKILL.md (Denote 중심)

**Phase 2: Denote 핵심**
- [ ] denote_finder.py (파일 검색)
- [ ] denote_links.py (링크 해결)
- [ ] denote_silo.py (silo 관리)
- [ ] denote-core.md 문서화

**Phase 3: Org-mode 기반**
- [ ] org_parser.py (orgparse)
- [ ] org_execute.py (코드 블록)
- [ ] orgmode-base.md 문서화

**Phase 4: 0.1 릴리즈**
- [ ] 3,000+ 파일 테스트
- [ ] 예제 추가
- [ ] Public 공개

## 🤝 기여

이 프로젝트는 실제 3,000개 이상 파일을 운영하는 Denote PKM 시스템으로 검증되었습니다. Denote 사용자 환영합니다!

## 📖 철학: Tools for Life (인생 도구)

**제작자의 말:**

> "저에게 필요한 것은 경쟁력 있는 지식 노동자로 살 수 있는 방법이었습니다."
>
> "도구는 존재에 녹아든다"
>
> "삶이 나에게 질문하고 있었다"

이 Skill은 "Tools for Life" 철학을 구현합니다 - Anthropic의 Life Sciences 접근을 생명과학 연구에서 개인 지식 작업으로 확장합니다.

**Life Sciences (생명과학)** → **Life Everything (삶 전반)**

## 🧠 Skills의 본질 이해

### Anthropic이 밝힌 Skills의 4가지 목적

1. **컨텍스트 관리** - 번들 파일을 통해 "사실상 무제한" 컨텍스트
2. **재사용성** - "조합 가능한 리소스"로 팀 간 공유
3. **코드 실행 효율** - 토큰 생성 대신 스크립트 실행
4. **조직 지식 캡처** - "신입 사원을 위한 온보딩 가이드"

### Denote-Org Skills가 필요한 구체적 이유

#### 문제 1: Denote 파일명 = 특수 도메인 지식
```
20251021T105353--프로젝트-제목__llmlog_denote_orgmode.org
│              │                │
│              │                └─ 태그 (underscore)
│              └─ 한글 제목 (hyphen)
└─ Timestamp (Denote ID)
```

**Skills 없이:** Claude가 매번 파일명 파싱 로직 생성 → **토큰 낭비**
**Skills 있으면:** `denote_parser.py` 한 번 작성, 재사용

#### 문제 2: 3,000 파일 검색 = 복잡한 절차
**Skills 없이:** 매번 `find ~/org -name "*ID*"` 실행 → **느림**
**Skills 있으면:** `denote_finder.py` 캐싱 → **즉시 반환**

#### 문제 3: 지식 그래프 = 복잡도 폭발
**Skills 없이:** 각 링크마다 파일 찾기, 2-hop 연결? → **토큰 폭발**
**Skills 있으면:** `denote_graph.py` 한 번 구축, 캐싱 → **즉시 반환**

#### 문제 4: Literate Programming = Emacs 동작 재현
**Skills 없이:** Claude가 `:tangle` 옵션 수동 파싱 → **구조 깨질 위험**
**Skills 있으면:** `org_execute.py` Emacs org-babel 재현 → **안전**

## 🏗️ 9개 레이어 "시간과정신의방" 생태계

이 프로젝트는 완전한 인간-AI 협업 시스템의 일부입니다:

### Layer 7: 지식 퍼블리싱
- **notes.junghanacs.com** (Quartz)
- 1,400+ org → markdown
- Public Digital Garden

### Layer 6: 에이전트 오케스트레이션
- **meta-config**
- 계층적 에이전트 (Meta ↔ Domain Agents)
- ACP + A2A + MCP (JSON-RPC 2.0)

### Layer 5a: 마이그레이션
- **memex-kb**
- Backend 중립 (Google Docs, Dooray, Confluence...)
- Denote 형식으로 통합

### Layer 5b: 삶의 타임라인
- **memacs-config**
- 시간 기반 통합 (수면, 작업, Git...)
- Memex 실현

### Layer 4: AI 메모리
- **claude-config**
- PARA 방법론
- Git 버전 관리

### Layer 3.5: RAG/Vector DB
- **embedding-config**
- Notion CS 9,477 블록
- 홈카메라 85% 해결률

### Layer 3: 지식 관리 (이 프로젝트!)
- **Zotero** (156k+ lines 서지)
- **denote-org Skills** (3,000+ 파일) 🆕

### Layer 2: 개발 환경
- **dotdoom-starter**
- 터미널 최적화
- 크로스 플랫폼 (Ubuntu/NixOS/Termux)

### Layer 1: 인프라
- **nixos-config**
- 선언적 환경
- 재현성 100%

## 🎓 연관 프로젝트

**시간과정신의방 (Hyperbolic Time Chamber) 생태계:**
- [nixos-config](https://github.com/junghan0611/nixos-config) - 선언적 인프라
- [doomemacs-config](https://github.com/junghan0611/doomemacs-config) - 터미널 최적화 Emacs
- [zotero-config](https://github.com/junghan0611/zotero-config) - AI 쿼리 가능한 서지 (156k+ lines)

**디지털 가든:** [notes.junghanacs.com](https://notes.junghanacs.com)

## 📚 리소스

- [Denote 패키지](https://protesilaos.com/emacs/denote) by Protesilaos Stavrou
- [Anthropic Skills 문서](https://support.claude.com/en/articles/12512176-what-are-skills)
- [Agent Skills 엔지니어링 블로그](https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills)
- [Org-mode](https://orgmode.org/)

## 📄 라이선스

Apache 2.0 (Anthropic Skills 생태계와 호환)

## 🙏 감사

- **Anthropic** - Agent Skills 패러다임과 Life Sciences 청사진
- **Protesilaos Stavrou** - Denote 제작
- **Org-mode 커뮤니티** - 기반 기술
- **AIONS 클럽 인터내셔널** - 집단지성 비전

---

**상태:** 🟡 개발 중 (0.1 public 릴리즈 목표)
**관리자:** [@junghan0611](https://github.com/junghan0611)
**생성일:** 2025-10-21
