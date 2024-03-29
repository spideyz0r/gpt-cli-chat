%global go_version 1.18.10
%global go_release go1.18.10

Name:           gpt-cli-chat
Version:        0.5
Release:        1%{?dist}
Summary:        gpt-cli-chat tool
License:        GPLv3
URL:            https://github.com/spideyz0r/gpt-cli-chat
Source0:        %{url}/archive/refs/tags/v%{version}.tar.gz

BuildRequires:  golang >= %{go_version}
BuildRequires:  git

%description
gpt-cli-chat is a cli tool to use open ai chat models.

%global debug_package %{nil}

%prep
%autosetup -n %{name}-%{version}

%build
go build -v -o %{name} -ldflags=-linkmode=external

%check
go test

%install
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}

%files
%{_bindir}/gpt-cli-chat

%license LICENSE

%changelog
* Sun Nov 19 2023 spideyz0r <47341410+spideyz0r@users.noreply.github.com> 0.5-1
- Add support to Internet queries

* Tue Aug 15 2023 spideyz0r <47341410+spideyz0r@users.noreply.github.com> 0.2-1
- Add support to GPT 4

* Thu Apr 13 2023 spideyz0r <47341410+spideyz0r@users.noreply.github.com> 0.2-1
- Bump version

* Thu Apr 13 2023 spideyz0r <47341410+spideyz0r@users.noreply.github.com> 0.1-1
- Initial build

