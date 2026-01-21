import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:intl/intl.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../models/subscription.dart';

class SubscriptionService {
  static const String baseUrl = 'http://10.0.2.2:5000';

  /// 🔐 Get auth headers with JWT
  Future<Map<String, String>> _authHeaders() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('token');

    if (token == null) {
      throw Exception('Not authenticated');
    }

    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    };
  }

  /// 📄 Fetch subscriptions (JWT protected)
  Future<List<Subscription>> fetchSubscriptions() async {
    final headers = await _authHeaders();

    final response = await http.get(
      Uri.parse('$baseUrl/api/subscriptions'),
      headers: headers,
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to load subscriptions');
    }

    final data = jsonDecode(response.body);
    final List listData = (data as List? ?? []);
    return listData.map((e) => Subscription.fromJson(e)).toList();
  }

  /// ➕ Create subscription (NO user_id)
  Future<void> createSubscription({
    required String name,
    required String category,
    required double amount,
    required String currency,
    required String billingCycle,
    required DateTime billingDate,
    required String status,
    required bool isTrial,
  }) async {
    final headers = await _authHeaders();

    final response = await http.post(
      Uri.parse('$baseUrl/api/subscriptions'),
      headers: headers,
      body: jsonEncode({
        'name': name,
        'category': category,
        'amount': amount,
        'currency': currency,
        'billing_cycle': billingCycle,
        'billing_date': DateFormat('yyyy-MM-dd').format(billingDate),
        'status': status,
        'is_trial': isTrial,
      }),
    );

    if (response.statusCode != 201) {
      throw Exception(response.body);
    }
  }

  /// 🗑 Delete subscription (JWT protected)
  Future<void> deleteSubscription(int id) async {
    final headers = await _authHeaders();

    final response = await http.delete(
      Uri.parse('$baseUrl/api/subscriptions/$id'),
      headers: headers,
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to delete subscription');
    }
  }
}
