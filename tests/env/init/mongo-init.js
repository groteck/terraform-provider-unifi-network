db = db.getSiblingDB('unifi');

db.createUser({
  user: 'unifi',
  pwd: 'unifi',
  roles: [
    { role: 'readWrite', db: 'unifi' },
    { role: 'dbAdmin', db: 'unifi' }
  ]
});

db.createCollection('setting');
db.setting.insertOne({
  "key": "is_setup",
  "value": true
});

// Inject a default admin user (admin/password123)
db.admin.insertOne({
  "name": "admin",
  "x_shadow": "$6$967AE4B000000000$8ED992C00000000008ED992C000000000967AE4B0000000067EC5CD0:1577279704",
  "email": "admin@local.host",
  "last_site_name": "default"
});

// Inject the API Key for token testing
db.apikey.insertOne({
  "name": "terraform-test",
  "key": "tf-test-token-12345",
  "site_id": "default",
  "permissions": ["all"]
});
