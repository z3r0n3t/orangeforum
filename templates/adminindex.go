package templates

const adminindexSrc = `
{{ define "content" }}

Yo
{{ index .Config "ForumName" }}

<form action="/admin" method="POST">

<input type="hidden" name="csrf" value="{{ .CSRF }}"><br>
Forum Name: <input type="text" name="forum_name" value="{{ index .Config "forum_name" }}" required><br>
Announcement: <input type="text" name="header_msg" value="{{ index .Config "header_msg" }}"><br>
<input type="checkbox" name="signup_disabled" value="1"{{ if index .Config "signup_disabled" }} checked{{ end }}>Signup disabled<br>
<input type="checkbox" name="group_creation_disabled" value="1"{{ if index .Config "group_creation_disabled" }} checked{{ end }}>Group creation disabled<br>
<input type="checkbox" name="image_upload_enabled" value="1"{{ if index .Config "image_upload_enabled" }} checked{{ end }}>Allow image upload<br>
<input type="checkbox" name="file_upload_enabled" value="1"{{ if index .Config "file_upload_enabled" }} checked{{ end }}>Allow file upload<br>
<input type="checkbox" name="allow_group_subscription" value="1"{{ if index .Config "allow_group_subscription" }} checked{{ end }}>Allow e-mail subscriptions to groups<br>
<input type="checkbox" name="allow_topic_subscription" value="1"{{ if index .Config "allow_topic_subscription" }} checked{{ end }}>Allow e-mail subscriptions to topics<br>
Data Directory: <input type="text" name="data_dir" value="{{ index .Config "data_dir" }}"><br>
FROM E-mail: <input type="text" name="default_from_mail" value="{{ index .Config "default_from_mail" }}"><br>
SMTP Host: <input type="text" name="smtp_host" value="{{ index .Config "smtp_host" }}"><br>
SMTP Port: <input type="number" name="smtp_port" value="{{ index .Config "smtp_port" }}"><br>
SMTP Username: <input type="text" name="smtp_user" value="{{ index .Config "smtp_user" }}"><br>
SMTP Password: <input type="text" name="smtp_pass" value="{{ index .Config "smtp_pass" }}"><br>

<input type="submit" value="Update">

</form>

{{ .Msg }}

Number of users: {{ .NumUsers }}<br>
Number of groups: {{ .NumGroups }}<br>
Number of topics: {{ .NumTopics }}<br>
Number of comments: {{ .NumComments }}<br>

{{ end }}`