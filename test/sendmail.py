#!/usr/bin/env python3
import argparse
import smtplib
import ssl
import sys
from email.message import EmailMessage
from pathlib import Path

def parse_args():
    parser = argparse.ArgumentParser(description="Send an email via SMTP with STARTTLS, support for self-signed certificates.")
    parser.add_argument("--server", required=True, help="SMTP server address")
    parser.add_argument("--port", type=int, default=587, help="SMTP server port (default: 587)")
    parser.add_argument("--username", required=False, help="SMTP username")
    parser.add_argument("--password", required=False, help="SMTP password")
    parser.add_argument("--from", dest="from_addr", required=True, help="From email address")
    parser.add_argument("--to", dest="to_addrs", required=True, nargs="+", help="Recipient email addresses")
    parser.add_argument("--subject", default="(no subject)", help="Email subject")
    parser.add_argument("--body-file", required=True, help="Path to the file containing the email body")
    parser.add_argument("--html", action="store_true", help="Send email as HTML")
    parser.add_argument("--insecure", action="store_true", help="Ignore TLS certificate verification (use with self-signed certs)")
    return parser.parse_args()

def main():
    args = parse_args()

    body_path = Path(args.body_file)
    if not body_path.exists():
        print(f"Error: Body file not found: {body_path}", file=sys.stderr)
        sys.exit(1)
    body = body_path.read_text(encoding="utf-8")

    msg = EmailMessage()
    msg["From"] = args.from_addr
    msg["To"] = ", ".join(args.to_addrs)
    msg["Subject"] = args.subject
    subtype = "html" if args.html else "plain"
    msg.set_content(body, subtype=subtype)

    # ignore ssl certificate verification if insecure flag is set
    if args.insecure:
        context = ssl.create_default_context()
        context.check_hostname = False
        context.verify_mode = ssl.CERT_NONE
    else:
        context = ssl.create_default_context()

    try:
        with smtplib.SMTP(args.server, args.port) as smtp:
            smtp.ehlo()
            smtp.starttls(context=context)
            smtp.ehlo()
            if args.username and args.password:
               smtp.login(args.username, args.password)
            smtp.send_message(msg)
        print("send success")
    except Exception as e:
        print(f"failed to send: {e}", file=sys.stderr)
        sys.exit(2)

if __name__ == "__main__":
    main()

