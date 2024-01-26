#
# spec file for package suse-uptime-tracker
#
# Copyright (c) 2024 SUSE LLC
#
# All modifications and additions to the file contributed by third parties
# remain the property of their copyright owners, unless otherwise agreed
# upon. The license for this file, and modifications and additions to the
# file, is the same license as for the pristine package itself (unless the
# license for the pristine package is not an Open Source License, in which
# case the license is the MIT License). An "Open Source License" is a
# license that conforms to the Open Source Definition (Version 1.9)
# published by the Open Source Initiative.

# Please submit bugfixes or comments via https://bugs.opensuse.org/
#


%global provider_prefix github.com/SUSE/uptime-tracker
%global import_path     %{provider_prefix}

Name:           suse-uptime-tracker
# the version will get set by the 'set_version' service
Version:        1.0.0
Release:        0
URL:            https://github.com/SUSE/uptime-tracker
License:        LGPL-2.1-or-later
Summary:        Service to track system uptime
Group:          System/Management
Source:         uptime-tracker-%{version}.tar.xz
Source1:        %name-rpmlintrc
BuildRequires:  golang-packaging
BuildRequires:  go1.18-openssl

ExcludeArch:    %ix86 s390 ppc64

%description
This package provides a utility service to track system uptime.

%{go_provides}

%prep
%setup -q -n uptime-tracker-%{version}

%build
find %_builddir/..
echo %{version} > suse-uptime-tracker/version.txt
%goprep %{import_path}
find %_builddir/..
go list all
%gobuild suse-uptime-tracker
find %_builddir/..

%install
%goinstall

# Install the suse-uptime-tracker timer and service.
install -D -m 644 %_builddir/go/src/%import_path/suse-uptime-tracker.timer %buildroot/%_unitdir/suse-uptime-tracker.timer
install -D -m 644 %_builddir/go/src/%import_path/suse-uptime-tracker.service %buildroot/%_unitdir/suse-uptime-tracker.service
ln -sf service %buildroot/%_sbindir/rcsuse-uptime-tracker

find %_builddir/..
# we currently do not ship the source for any go module
rm -rf %buildroot/usr/share/go

%pre
%service_add_pre suse-uptime-tracker.service suse-uptime-tracker.timer

%post
# SLES12 systemd does not support RandomizedDelaySec so remove it
%if (0%{?sle_version} > 0 && 0%{?sle_version} < 150000)
    sed -i '/RandomizedDelaySec*/d' %{_unitdir}/suse-uptime-tracker.timer
%endif
%service_add_post suse-uptime-tracker.service suse-uptime-tracker.timer

%preun
%service_del_preun suse-uptime-tracker.service suse-uptime-tracker.timer

%postun
%service_del_postun suse-uptime-tracker.service suse-uptime-tracker.timer
if [ -e /etc/zypp/suse-uptime.logs ]; then
  rm -f /etc/zypp/suse-uptime.logs 2> /dev/null
fi

%check
%gotest -v %import_path/suse-uptime-tracker
make -C %_builddir/go/src/%import_path gofmt

%files
%license LICENSE LICENSE.LGPL
%doc README.md
%_bindir/suse-uptime-tracker
%_sbindir/rcsuse-uptime-tracker
%_unitdir/suse-uptime-tracker.service
%_unitdir/suse-uptime-tracker.timer

%changelog
