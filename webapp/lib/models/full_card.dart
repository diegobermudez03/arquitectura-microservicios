class FullCard {
  // User fields
  final String userId;
  final String userToken;
  final String userName;
  final String userLastname;
  final String userBirthDate;
  final String userCountryCode;
  final String userSocialId;
  final String userCreatedAt;

  // Card fields
  final String cardId;
  final String cardPan;
  final String cardCvv;
  final String cardExpiry;
  final String cardType;
  final String cardStatus;
  final String cardCreatedAt;

  FullCard({
    required this.userId,
    required this.userToken,
    required this.userName,
    required this.userLastname,
    required this.userBirthDate,
    required this.userCountryCode,
    required this.userSocialId,
    required this.userCreatedAt,
    required this.cardId,
    required this.cardPan,
    required this.cardCvv,
    required this.cardExpiry,
    required this.cardType,
    required this.cardStatus,
    required this.cardCreatedAt,
  });

  factory FullCard.fromJson(Map<String, dynamic> json) => FullCard(
        userId: json['user_id'] ?? '',
        userToken: json['user_token'] ?? '',
        userName: json['user_name'] ?? '',
        userLastname: json['user_lastname'] ?? '',
        userBirthDate: json['user_birth_date'] ?? '',
        userCountryCode: json['user_country_code'] ?? '',
        userSocialId: json['user_social_id'] ?? '',
        userCreatedAt: json['user_created_at'] ?? '',
        cardId: json['card_id'] ?? '',
        cardPan: json['card_pan'] ?? '',
        cardCvv: json['card_cvv'] ?? '',
        cardExpiry: json['card_expiry'] ?? '',
        cardType: json['card_type'] ?? '',
        cardStatus: json['card_status'] ?? '',
        cardCreatedAt: json['card_created_at'] ?? '',
      );
}
