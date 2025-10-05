import 'package:flutter_dotenv/flutter_dotenv.dart';

class Env {
  static String get issueServiceUrl {
    final url = 'https://issuer_url_here.com';
    print('Issue Service URL: $url');
    return url;
  }
  
  static String get notificationServiceUrl {
    final url = 'https://notifications_url_here.com';
    print('Notification Service URL: $url');
    return url;
  }
}
