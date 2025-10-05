class DeclineReason {
  final String reason;
  
  DeclineReason({required this.reason});
  
  factory DeclineReason.fromJson(Map<String, dynamic> json) =>
      DeclineReason(reason: json['reason']);
}

class IssuedCard {
  final String pan;
  final String cvv;
  final String expiryDate;
  final String cardType;
  
  IssuedCard({
    required this.pan,
    required this.cvv,
    required this.expiryDate,
    required this.cardType,
  });
  
  factory IssuedCard.fromJson(Map<String, dynamic> json) => IssuedCard(
      pan: json['pan'],
      cvv: json['cvv'],
      expiryDate: json['expiry_date'],
      cardType: json['card_type']);
}

class IssuerResponse {
  final DeclineReason? declineReason;
  final IssuedCard? issuedCard;
  final String requestUUID;
  final String suscriptorToken;
  final String status;

  IssuerResponse({
    this.declineReason,
    this.issuedCard,
    required this.requestUUID,
    required this.suscriptorToken,
    required this.status,
  });

  factory IssuerResponse.fromJson(Map<String, dynamic> json) => IssuerResponse(
    declineReason: json['decline_reason'] != null 
        ? DeclineReason.fromJson(json['decline_reason']) 
        : null,
    issuedCard: json['issued_card'] != null 
        ? IssuedCard.fromJson(json['issued_card']) 
        : null,
    requestUUID: json['request_uuid'],
    suscriptorToken: json['suscriptor_token'],
    status: json['status'],
  );
}
