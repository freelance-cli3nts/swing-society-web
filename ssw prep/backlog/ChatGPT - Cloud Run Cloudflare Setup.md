# Cloud Run Cloudflare Setup

1. Set Up Your Cloud Run Service

Deploy your app to Cloud Run:
Bashgcloud run deploy SERVICE_NAME --image IMAGE_URL --platform managed

Get the URL of your Cloud Run service:
Bashgcloud run services describe SERVICE_NAME --platform managed --region REGION --format"value(status.url)"

2. Configure Custom Domain in Cloud Run

Add your domain:
Bashgcloud run domain-mappings create --service SERVICE_NAME --domain YOUR_DOMAIN

Verify the domain:
Google will give you DNS records (CNAMEs) to set up in Cloudflare.

3. Set Up DNS on Cloudflare

Go to Cloudflare's DNS tab.

Add a CNAME record:

Name: www or yourdomain.com

Target: The URL Google gave you (like ghs.googlehosted.com)

Proxy status: Set it to "DNS only" for now (orange cloud off).

Add TXT records (if Google asks for them to verify ownership).

4. Configure SSL (Optional but Recommended)

By default, Cloud Run serves traffic over HTTPS. Cloudflare has different SSL modes:

Full (strict): Ideal — needs a valid certificate from Google (which Cloud Run provides).

Full: Works, but less secure.

Flexible: Avoid — causes redirect loops with Cloud Run.

Go to SSL/TLS in Cloudflare and set it to Full (strict) for the best security.

5. Test the Setup

Wait for DNS propagation (usually a few minutes to a couple of hours).

Visit your domain — it should now route to your Cloud Run app.

If it doesn’t work, clear Cloudflare cache or double-check your DNS records.

🔥 Bonus: If you want to enable Cloudflare’s CDN, caching, or performance features, turn the orange cloud back on (proxy mode) after everything works.

Want me to create a script or Terraform setup to automate this? 🚀

### Troubleshooting checklist:

1. Verify DNS Propagation
Ensure your DNS records are live:
Use a DNS checker like dnschecker.org to confirm your CNAME and TXT records are propagated globally.
For a quick terminal check:
Bash nslookup yourdomain.comdig yourdomain.com

If they’re missing or incorrect, double-check the Cloudflare DNS setup.

1. Check Cloud Run Domain Mapping Status

Let’s make sure Google Cloud recognized the domain:

```
Bashgcloud run domain-mappings describe --domain yourdomain.com
```

It should show ACTIVE.
If not, try redeploying the mapping:

```
Bashgcloud run domain-mappings create --service YOUR_SERVICE_NAME --domain yourdomain.com
```

3. SSL/TLS Setup on Cloudflare

Ensure Cloudflare’s SSL/TLS mode is set to "Full (strict)".
If it’s Flexible, it causes redirect loops or 502 errors with Cloud Run.

4. Cloudflare Proxy (Orange Cloud)

If the orange cloud (proxy) is enabled, try turning it off (DNS only) temporarily to check if Cloudflare is causing the issue.
Once the domain works, you can turn the proxy back on for performance/caching benefits.

5. Firewall and Security Rules

Check that Cloudflare’s security settings aren’t blocking requests:

Security > Events — see if Cloudflare logs any blocked requests.

Make sure Google’s IP ranges aren’t blocked: Google IP Ranges

6. Check Cloud Run App Health

Ensure the app itself is reachable:

```
Bashcurl https://your-cloud-run-url
```

If that works but your domain doesn’t, the problem is DNS/Cloudflare-related.

Still stuck? If you share your DNS setup or the error you’re seeing, I can dig deeper. 🚀


7. Verify Domain Registration and Nameservers

First, ensure that swingsociety.eu is properly registered and that its nameservers are correctly set:

Domain Registration: Confirm that the domain registration is active and hasn't expired.

Nameserver Configuration: Check that the domain's nameservers are correctly configured to point to Cloudflare's nameservers. You can verify this through your domain registrar's control panel or by using a WHOIS lookup tool.

8. Inspect DNSSEC Settings

DNSSEC adds an extra layer of security to DNS but can cause SERVFAIL errors if misconfigured:

DNSSEC Status: Verify whether DNSSEC is enabled for your domain. If it is, ensure that all DNSSEC records are correctly set up. Misconfigurations here can lead to validation failures.

Disabling DNSSEC: If you're uncertain about the DNSSEC configuration or suspect issues, consider temporarily disabling it to see if the SERVFAIL


9. Verify DNS Records in Cloudflare

In Cloudflare:

Go to DNS settings for swingsociety.eu.

Add these A records:
Name: swingsociety.eu  Type: A  Content: 216.239.32.21  Name: swingsociety.eu  Type: A  Content: 216.239.34.21  Name: swingsociety.eu  Type: A  Content: 216.239.36.21  Name: swingsociety.eu  Type: A  Content: 216.239.38.21  

And these AAAA records (IPv6):
Name: swingsociety.eu  Type: AAAA  Content: 2001:4860:4802:32::15  Name: swingsociety.eu  Type: AAAA  Content: 2001:4860:4802:34::15  Name: swingsociety.eu  Type: AAAA  Content: 2001:4860:4802:36::15  Name: swingsociety.eu  Type: AAAA  Content: 2001:4860:4802:38::15  

💡 Important:

Set the Proxy status to DNS only (gray cloud). Cloudflare’s proxy (orange cloud) can interfere with Google’s certificate provisioning.

If you have a CNAME for www, ensure it points to your Cloud Run URL or ghs.googlehosted.com.

10. Check for DNS Propagation

After saving changes, confirm DNS is working globally:

```
Bashdig swingsociety.eu +short
```

You should see the IPs from Cloud Run (216.239.x.x).

Alternatively, use DNS Checker to ensure the world sees the records.

11. Trigger Certificate Retry

If DNS looks good, nudge Cloud Run to re-check the domain setup:

```
Bashgcloud run domain-mappings update swingsociety.eu --service ss-go
```

Or you can delete and recreate the mapping:

```
Bashgcloud run domain-mappings delete --domain swingsociety.eugcloud run domain-mappings create --service ss-go --domain swingsociety.eu
```

12. Test the Domain

Once the certificate is issued (it can take a few minutes), check the domain in your browser.
If it still fails, show me the updated status output or any errors, and I’ll guide you through the next steps! 
