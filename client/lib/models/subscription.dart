class Subscription {
  final int id;
  final String name;
  final String category;
  final double amount;
  final String currency;
  final String billingCycle;
  final DateTime billingDate;
  final String status;
  final bool isTrial;

  Subscription({
    required this.id,
    required this.name,
    required this.category,
    required this.amount,
    required this.currency,
    required this.billingCycle,
    required this.billingDate,
    required this.status,
    required this.isTrial,
  });

  factory Subscription.fromJson(Map<String, dynamic> json) {
    return Subscription(
      id: json['id'],
      name: json['name'],
      category: json['category'] ?? '',
      amount: (json['amount'] as num).toDouble(),
      currency: json['currency'],
      billingCycle: json['billing_cycle'],
      billingDate: DateTime.parse(json['billing_date']),
      status: json['status'],
      isTrial: json['is_trial'],
    );
  }
}
