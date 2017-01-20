%global import_path     github.com/cavaliercoder/redirector
%global commit          bd8f06e4337b27b289f26a60b5d28fd0e4ce7d2b                
%global shortcommit     %(c=%{commit}; echo ${c:0:7})

# https://bugzilla.redhat.com/show_bug.cgi?id=995136#c12
%global _dwz_low_mem_die_limit 0

Name:           redirector
Version:        1.1.2
Release:        1%{?dist}
Summary:        Simple HTTP redirect server

License:        MIT
URL:            https://%{import_path}
Source0:        %{name}-%{version}.tar.gz

Requires(pre):    /usr/sbin/useradd, /usr/bin/getent, /usr/bin/systemctl
Requires(postun): /usr/sbin/userdel, /usr/bin/systemctl

%description
Redirector is a fast and lightweight HTTP server that serves a single purpose;
to redirect web clients from one URL to another. This is useful in the following
situations:

* URL shortening
* migrating from one URL layout to another
* vanity URLs

Redirector uses BoltDB, an embedded, high performance key/value store to provide
sub-millisecond responses, even when managing millions of URL mappings.

%prep
%setup -q -n %{name}-%{version}
make get-deps

%build
make %{name}

%install
rm -rf %{buildroot}
install -d %{buildroot}%{_bindir}
install -d %{buildroot}%{_sysconfdir}/%{name}
install -d %{buildroot}%{_sharedstatedir}/%{name}
install -d %{buildroot}%{_localstatedir}/log/%{name}
install -d %{buildroot}%{_exec_prefix}/lib/systemd/system/
install -p -m 755 %{name} %{buildroot}%{_bindir}/%{name}

# create default configuration file
cat > %{buildroot}%{_sysconfdir}/%{name}/%{name}.json <<EOL
{
	"listenAddr": ":8080",
	"mgmtAddr": "127.0.0.1:9321",
	"logFile": "%{_localstatedir}/log/%{name}/%{name}.log",
	"accessLogFile": "%{_localstatedir}/log/%{name}/access.log",
	"database": "bolt",
	"databasePath": "%{_sharedstatedir}/%{name}/%{name}.db",
	"keyBuilder": "path"
}
EOL

# create systemd service
# TODO: service should run as system account
cat > %{buildroot}%{_exec_prefix}/lib/systemd/system/%{name}.service <<EOL
[Unit]
Description=Simple HTTP redirect server
After=network.target remote-fs.target nss-lookup.target

[Service]
Type=simple
User=%{name}
ExecStart=%{_bindir}/%{name} serve
KillSignal=SIGQUIT
TimeoutStopSec=5
KillMode=process

[Install]
WantedBy=multi-user.target
EOL

# TODO: add logrotated config

%pre
/usr/bin/getent group %{name} >/dev/null || /usr/sbin/groupadd -r %{name}
/usr/bin/getent passwd %{name} >/dev/null || /usr/sbin/useradd \
	--comment "%{name} web server" \
	--system \
	--home-dir %{_sharedstatedir}/%{name} \
	--shell /sbin/nologin \
	--gid %{name} \
	%{name}

%post
/usr/bin/systemctl daemon-reload || :

%preun
if [ "$1" = "0" ]; then
	/usr/bin/systemctl stop %{name} || :
fi

%postun
if [ "$1" = "0" ]; then
	/usr/sbin/userdel --force %{name} || :
	/usr/bin/systemctl daemon-reload || :
else
	/usr/bin/systemctl restart redirector
fi

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%dir %attr(0750, root, %{name}) %{_sysconfdir}/%{name}
%dir %attr(0750, %{name}, %{name}) %{_sharedstatedir}/%{name}
%dir %attr(0750, %{name}, %{name}) %{_localstatedir}/log/%{name}
%config(noreplace) %attr(0640, root, %{name}) %{_sysconfdir}/%{name}/%{name}.json
%attr(0644, root, root) %{_exec_prefix}/lib/systemd/system/%{name}.service
%attr(0755, root, root)%{_bindir}/%{name}

%changelog
* Fri Jan 20 2017 Ryan Armstrong <ryan@cavaliercoder.com> - 1.1.2-1
- Added more detail to access logging
- Added version to server start log entry

* Wed Jan 18 2017 Ryan Armstrong <ryan@cavaliercoder.com> - 1.1.1-1
- Improved import performance
- Improved error handling for key builders

* Thu Jan 12 2017 Ryan Armstrong <ryan@cavaliercoder.com> - 1.1.0-1
- Added destination templating
- Added server response header

* Wed Jan 11 2017 Ryan Armstrong <ryan@cavaliercoder.com> - 1.0.0-1
- Initial RPM release

